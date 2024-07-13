package shoptraininglogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/app/training/model/common"
	yqmongo "yufuture-gpt/app/training/model/mongo"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/consts"
	"yufuture-gpt/common/utils"
)

type PreSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type ApplyGoodsIdListReq struct {
	SerialNumber string        `json:"serialNumber"`
	Uuid         string        `json:"uuid"`
	ShopName     string        `json:"shop_name"`
	Platform     string        `json:"platform"`
	Token        string        `json:"token"`
	Cookies      []interface{} `json:"cookies"`
}

// FetchEstimateResultResp 预估训练所需消耗的资源
type FetchEstimateResultResp struct {
	Code int64 `json:"code"`
	Data struct {
		Token    int64 `json:"token"`
		Power    int64 `json:"power"`
		FileSize int64 `json:"filesize"`
	} `json:"data"`
}

type ApplyGoodsJsonReq struct {
	AppType string    `json:"app_type"`
	SkuList []*string `json:"sku_list"`
	Channel string    `json:"channel"`
}

type FetchGoodsJsonReq struct {
	AppType string `json:"app_type"`
	Channel string `json:"channel"`
	Limit   int    `json:"limit"`
}

type FetchGoodsJsonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Items []struct {
			Sku         string `json:"sku"`
			CollectTime string `json:"collectTime"`
			Url         string `json:"url"`
		} `json:"items"`
	} `json:"data"`
}

func NewPreSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingLogic {
	return &PreSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PreSettingLogic) PreSetting(in *training.ShopTrainingReq) (*training.ShopTrainingResp, error) {
	// 构建请求体 发送爬取商品ID请求
	serialNumber := l.svcCtx.SnowFlakeNode.Generate().String()
	err := ApplyGoodsListId(l.Logger, l.svcCtx, in.Uuid, in.ShopName, in.PlatformType, serialNumber, in.Authorization, in.Cookies)
	if err != nil {
		l.Logger.Error("发送爬取商品ID请求失败", err)
		return nil, err
	}
	// 等待两分钟
	time.Sleep(time.Minute * 2)
	// 轮询等待爬取商品ID落库完成 通过serialNumber在mongo中查找
	var saveShop *yqmongo.Shoptrainingshoptitles
	i := 0
	for i < 10 {
		saveShop, err = l.svcCtx.ShoptrainingshoptitlesModel.FindOneBySerialNumber(l.ctx, serialNumber)
		if err != nil {
			l.Logger.Error("根据serialNumber在mongo中查找店铺失败", err)
		} else if saveShop != nil {
			break
		}
		time.Sleep(time.Minute * 6)
		i++
	}
	// 没查到直接认为此次预设失败
	if saveShop == nil {
		l.Logger.Error("根据serialNumber在mongo中查找店铺失败", err)
		return nil, err
	}
	// 将店铺和商品置为预设状态
	// 根据uuid和userid从mysql中查找出店铺
	tsShop, err := l.svcCtx.TsShopModel.FindOneByUuidAndUserId(l.ctx, in.UserId, in.Uuid)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找店铺失败", err)
		return nil, err
	}
	// 根据店铺shopId从mysql中查找出enabled字段为2启用商品列表
	tsGoodsList, err := l.svcCtx.TsGoodsModel.FindEnabledListByShopId(l.ctx, tsShop.Id)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找商品失败", err)
		return nil, err
	}
	// 商品列表转换为map 便于后续查找
	var tsGoodsMap map[string]*orm.TsGoods
	if tsGoodsList != nil {
		tsGoodsMap = make(map[string]*orm.TsGoods)
		for _, tsGoods := range *tsGoodsList {
			tsGoodsMap[tsGoods.PlatformId] = tsGoods
		}
	}
	// 更新店铺状态为预设中
	err = UpdateShopPreSetting(l.ctx, l.svcCtx, tsShop, in.UserId)
	if err != nil {
		l.Logger.Error("修改店铺状态失败", in)
		return nil, err
	}
	// 更新商品状态为预设中 同时提取本次需要预设的商品列表
	var preSettingGoodsList []*orm.TsGoods
	for _, saveGoods := range saveShop.GoodsList {
		// 同时在mongo中并且enabled字段为2启用的商品即为本次需要预设的商品
		if tsGoods, ok := tsGoodsMap[saveGoods.PlatformId]; ok {
			// 排除掉状态已经在预设中/训练中/预设完成的商品
			if tsGoods.TrainingStatus == consts.Presetting || tsGoods.TrainingStatus == consts.Training || tsGoods.TrainingStatus == consts.PresettingComplete {
				continue
			}
			// 更新商品状态为预设中
			err = UpdateGoodsPreSetting(l.ctx, l.svcCtx, tsGoods, in.UserId)
			if err != nil {
				l.Logger.Error("修改商品状态失败", tsGoods)
				return nil, err
			}
			//将筛选出的商品添加到预设商品列表
			preSettingGoodsList = append(preSettingGoodsList, tsGoods)
		}
	}

	// 请求获取商品JSON
	err = ApplyGoodsJson(l.svcCtx, preSettingGoodsList)
	if err != nil {
		l.Logger.Info("发送获取商品JSON请求失败", err)
		return nil, err
	}
	// 等待2分钟
	time.Sleep(time.Minute * 2)
	// 每6分钟调用一次接口 连续10次失败则结束
	FetchAndSaveGoodsJson(l.Logger, l.ctx, l.svcCtx, preSettingGoodsList)
	// 从jSON解析的商品列表文档
	var goodsDocumentList []*common.PddGoodsDocument
	// 获取并解析商品JSON到结果文档列表
	GetAndParseGoodsJson(l.Logger, tsShop, goodsDocumentList, preSettingGoodsList)
	// 构建获取训练时长的请求图片列表
	var goodPicList []string
	for _, goodsDocument := range goodsDocumentList {
		goodPicList = append(goodPicList, goodsDocument.PictureUrlList...)
	}

	// 发送请求 获取商品训练所需资源和时长
	var fetchEstimateResultResp FetchEstimateResultResp
	err = utils.HTTPPostAndParseJSON(l.svcCtx.Config.TrainingGoodsConf.FetchEstimateResultUrl, struct {
		Urls []string `json:"urls"`
	}{Urls: goodPicList}, &fetchEstimateResultResp)
	if err != nil {
		l.Logger.Error("获取商品训练所需资源和时长失败", err)
	}

	// 设计结构化文档 预设结果保存到mongo 正式训练时直接从mongo中取
	shoppresettingshoptitles := &yqmongo.Shoppresettingshoptitles{
		ShopId:     tsShop.Id,
		PlatformId: goodsDocumentList[0].PlatformMallId,
		UUID:       in.Uuid,
		UserID:     in.UserId,

		PreSettingToken:    fetchEstimateResultResp.Data.Token,
		PresettingPower:    fetchEstimateResultResp.Data.Power,
		PresettingFileSize: fetchEstimateResultResp.Data.FileSize,
		//PreSettingTime:
		GoodsDocumentList: goodsDocumentList,
	}
	err = l.svcCtx.ShoppresettingshoptitlesModel.Insert(l.ctx, shoppresettingshoptitles)
	if err != nil {
		l.Logger.Error("保存训练到mongo失败", err)
		return nil, err
	}
	// 更新店铺和商品状态为预设完成
	err = UpdateShopPreSettingComplete(l.ctx, l.svcCtx, tsShop, in.UserId)
	if err != nil {
		l.Logger.Error("修改店铺状态失败", err)
	}
	for _, preSettingGoods := range preSettingGoodsList {
		// 更新商品状态为预设完成
		err = UpdateGoodsPreSettingComplete(l.ctx, l.svcCtx, preSettingGoods, in.UserId)
		if err != nil {
			l.Logger.Error("修改商品状态失败", err)
		}
	}
	return &training.ShopTrainingResp{}, nil
}

func UpdateShopPreSetting(ctx context.Context, svcCtx *svc.ServiceContext, tsShop *orm.TsShop, userId int64) error {
	tsShop.TrainingStatus = consts.Presetting
	tsShop.UpdateTime = time.Now()
	tsShop.UpdateBy = userId
	err := svcCtx.TsShopModel.Update(ctx, tsShop)
	if err != nil {
		return err
	}
	return nil
}

func UpdateShopPreSettingComplete(ctx context.Context, svcCtx *svc.ServiceContext, tsShop *orm.TsShop, userId int64) error {
	tsShop.TrainingStatus = consts.PresettingComplete
	tsShop.UpdateTime = time.Now()
	tsShop.UpdateBy = userId
	err := svcCtx.TsShopModel.Update(ctx, tsShop)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGoodsPreSetting(ctx context.Context, svcCtx *svc.ServiceContext, tsGoods *orm.TsGoods, userId int64) error {
	tsGoods.TrainingStatus = consts.Presetting
	tsGoods.UpdateTime = time.Now()
	tsGoods.UpdateBy = userId
	err := svcCtx.TsGoodsModel.Update(ctx, tsGoods)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGoodsPreSettingComplete(ctx context.Context, svcCtx *svc.ServiceContext, tsGoods *orm.TsGoods, userId int64) error {
	tsGoods.TrainingStatus = consts.PresettingComplete
	tsGoods.UpdateTime = time.Now()
	tsGoods.UpdateBy = userId
	err := svcCtx.TsGoodsModel.Update(ctx, tsGoods)
	if err != nil {
		return err
	}
	return nil
}

func ApplyGoodsListId(log logx.Logger, svcCtx *svc.ServiceContext, uuid string, shopName string, platformType int64, serialNumber string, token string, cookie string) error {
	var platform string
	if platformType == consts.Pdd {
		platform = "pdd"
	}
	var cookies []interface{}
	err := utils.StringToAny(cookie, &cookies)
	if err != nil {
		log.Error("cookie转换cookies数组失败", err)
		return err
	}
	applyGoodsIdReq := &ApplyGoodsIdListReq{
		SerialNumber: serialNumber,
		Uuid:         uuid,
		ShopName:     shopName,
		Platform:     platform,
		Token:        token,
		Cookies:      []interface{}{},
	}
	var applyGoodsIdResp interface{}
	err = utils.HTTPPostAndParseJSON(svcCtx.Config.TrainingGoodsConf.ApplyGoodsIdListUrl, applyGoodsIdReq, &applyGoodsIdResp)
	if err != nil {
		log.Error("发送获取商品ID列表请求失败", err)
		return err
	}
	return nil
}

func ApplyGoodsJson(svcCtx *svc.ServiceContext, trainingGoodsList []*orm.TsGoods) error {
	var skuList []*string
	for _, trainingGoods := range trainingGoodsList {
		skuList = append(skuList, &trainingGoods.PlatformId)
	}
	// 发送获取商品JSON请求
	var applyGoodsJsonResp interface{}
	err := utils.HTTPPostAndParseJSON(svcCtx.Config.TrainingGoodsConf.ApplyGoodsJsonUrl, ApplyGoodsJsonReq{
		AppType: "pdd",
		Channel: svcCtx.Config.TrainingGoodsConf.ApplyGoodsJsonChannel,
		SkuList: skuList,
	}, &applyGoodsJsonResp)
	if err != nil {
		return err
	}
	return nil
}

// FetchAndSaveGoodsJson 获取获取商品json的url并更新到数据库
func FetchAndSaveGoodsJson(log logx.Logger, ctx context.Context, svcCtx *svc.ServiceContext, trainingGoodsList []*orm.TsGoods) {
	i := 0
	for i < 10 {
		var fetchGoodsJsonResp FetchGoodsJsonResp
		// 调用接口
		err := utils.HTTPPostAndParseJSON(svcCtx.Config.TrainingGoodsConf.FetchGoodsJsonUrl, FetchGoodsJsonReq{
			AppType: "pdd",
			Channel: svcCtx.Config.TrainingGoodsConf.ApplyGoodsJsonChannel,
			Limit:   100,
		}, &fetchGoodsJsonResp)
		// 将返回的获取商品json的url更新进mysql
		for _, item := range fetchGoodsJsonResp.Data.Items {
			err = svcCtx.TsGoodsModel.UpdateGoodsJsonUrlByPlatformId(ctx, item.Sku, item.Url)
			if err != nil {
				log.Error("更新商品json的url失败", err)
			}
		}
		// 查询mysql判断是否完成 同时将完成的新获取json的url保存进training_url
		complete := true
		for _, trainingGoods := range trainingGoodsList {
			// 已经更新获取json的url的商品跳过循环
			if trainingGoods.GoodsUrl != "" {
				continue
			}
			mysqlGoods, err := svcCtx.TsGoodsModel.FindOne(ctx, trainingGoods.Id)
			if err != nil {
				log.Error("查询训练中的商品数据失败", err)
			} else if mysqlGoods.GoodsJsonUrl == "" {
				complete = false
			} else {
				trainingGoods.GoodsUrl = mysqlGoods.GoodsJsonUrl
			}
		}
		if complete {
			//都完成后跳出循环
			break
		}
		time.Sleep(time.Minute * 6)
		i++
	}
}

func GetAndParseGoodsJson(log logx.Logger, tsShop *orm.TsShop, goodsDocumentList []*common.PddGoodsDocument, trainingGoodsList []*orm.TsGoods) {
	for _, trainingGoods := range trainingGoodsList {
		//根据获取商品JSON的url获取商品JSON 排除掉url为空的商品
		if trainingGoods.GoodsJsonUrl == "" {
			continue
		}
		var goodsJson string
		err := utils.HTTPGetAndParseJSON(trainingGoods.GoodsUrl, &goodsJson)
		if err != nil {
			log.Error("根据url获取商品json数据失败", err)
			continue
		}
		if goodsJson == "" {
			log.Error("根据url获取商品json数据失败, 返回的json为空", err)
			continue
		}

		// 解析JSON 将图片列表等数据保存下来
		goodsDocument := ParsePddGoods(goodsJson, tsShop, trainingGoods)
		goodsDocumentList = append(goodsDocumentList, goodsDocument)
	}
}
