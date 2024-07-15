package shoptraininglogic

import (
	"context"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/thirdparty"
	"yufuture-gpt/app/training/model/common"
	yqmongo "yufuture-gpt/app/training/model/mongo"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/consts"
	"yufuture-gpt/common/utils"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSettingShopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
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

func NewPreSettingShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingShopLogic {
	return &PreSettingShopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PreSettingShop 预设店铺
func (l *PreSettingShopLogic) PreSettingShop(in *training.PreSettingShopReq) (*training.PreSettingShopResp, error) {
	// 构建请求体 发送爬取商品ID请求
	serialNumber := l.svcCtx.SnowFlakeNode.Generate().String()
	err := thirdparty.ApplyGoodsListId(l.Logger, l.svcCtx, in.Uuid, in.ShopName, in.PlatformType, serialNumber, in.Authorization, in.Cookies)
	if err != nil {
		l.Logger.Error("发送爬取商品ID请求失败", err)
		return nil, err
	}
	// 等待两分钟
	time.Sleep(time.Minute * 2)
	// 轮询等待爬取商品ID落库完成 通过serialNumber在mongo中查找
	var saveShop *yqmongo.Dbsavegoodscrawlertitles
	i := 0
	for i < 10 {
		saveShop, err = l.svcCtx.DbsavegoodscrawlertitlesModel.FindOneBySerialNumber(l.ctx, serialNumber)
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
	err = thirdparty.ApplyGoodsJson(l.svcCtx, preSettingGoodsList)
	if err != nil {
		l.Logger.Info("发送获取商品JSON请求失败", err)
		return nil, err
	}
	// 等待2分钟
	time.Sleep(time.Minute * 2)
	// 每6分钟调用一次接口 连续10次失败则结束
	thirdparty.FetchAndSaveGoodsJson(l.Logger, l.ctx, l.svcCtx, preSettingGoodsList)
	// 从jSON解析的商品列表文档
	var goodsDocumentList []*common.PddGoodsDocument
	// 获取并解析商品JSON到结果文档列表
	thirdparty.GetAndParseGoodsJson(l.Logger, tsShop, goodsDocumentList, preSettingGoodsList)
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
	dbpresettingshoptitlesModel := &yqmongo.Dbpresettingshoptitles{
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
	err = l.svcCtx.DbpresettingshoptitlesModel.Insert(l.ctx, dbpresettingshoptitlesModel)
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
	return &training.PreSettingShopResp{}, nil
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
