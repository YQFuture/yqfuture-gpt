package shoptraininglogic

import (
	"context"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/utils"
)

type TrainingShopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// PddGoodsDocument 保存到es的拼多多商品训练结果
type PddGoodsDocument struct {
	ShopId  int64  `json:"shopId"`  // 对应mysql店铺id
	GoodsId int64  `json:"goodsId"` // 对应mysql商品id
	Uuid    string `json:"uuid"`
	UserId  int64  `json:"userId"`

	PlatformMallId  string `json:"platformMallId"`  // 对应json店铺id
	PlatformGoodsId string `json:"platformGoodsId"` // 对应json商品id
	GoodsUrl        string `json:"goodsUrl"`
	GoodsJson       string `json:"goodsJson"`
	GoodsName       string `json:"goodsName"`

	SkuSpecs                 map[string]string `json:"skuSpecs"`
	GroupPrice               float64           `json:"groupPrice"`
	NormalPrice              float64           `json:"normalPrice"`
	ServicePromise           map[string]string `json:"servicePromise"`           // 商品服务承诺列表
	SellPointTagList         interface{}       `json:"sellPointTagList"`         // 卖点
	PromptExplain            string            `json:"promptExplain"`            // 商品提示
	DetailGalleryDescription string            `json:"detailGalleryDescription"` // 图片训练结果描述

	PictureUrlList []string  `json:"pictureUrlList"`
	Token          int64     `json:"token"` // 消耗的token
	CreatedAt      time.Time `json:"createdAt"`
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

type ImageInfo struct {
	ID   string   `json:"id"`
	URLs []string `json:"urls"`
}

type CreateBatchTaskReq struct {
	SystemPrompt string       `json:"system_prompt"`
	BatchImages  []*ImageInfo `json:"batch_image_urls"`
}

func NewTrainingShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrainingShopLogic {
	return &TrainingShopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 训练店铺
func (l *TrainingShopLogic) TrainingShop(in *training.TrainingShopReq) (*training.TrainingShopResp, error) {
	// 根据uuid和userId从mongo中找到最新的一条店铺数据
	saveShop, err := l.svcCtx.ShoptrainingshoptitlesModel.FindNewOneByUuidAndUserId(l.ctx, in.Uuid, in.UserId)
	if err != nil {
		l.Logger.Error("根据uuid和userId从mongo中找到最新的一条店铺数据失败", err)
		return nil, err
	}

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

	// 更新店铺状态为训练中 添加训练次数
	err = UpdateShopTraining(l.ctx, l.svcCtx, tsShop, in.UserId)
	if err != nil {
		l.Logger.Error("修改店铺状态失败", in)
		return nil, err
	}
	// 更新商品状态为训练中 同时提取本次需要训练的商品列表
	var trainingGoodsList []*orm.TsGoods
	for _, saveGoods := range saveShop.GoodsList {
		// 同时在mongo中并且enabled字段为2启用的商品即为本次需要训练的商品
		if tsGoods, ok := tsGoodsMap[saveGoods.PlatformId]; ok {
			// 排除掉已经在训练中的商品
			if tsGoods.TrainingStatus == 1 {
				continue
			}
			// 更新商品状态为训练中 添加训练次数
			err = UpdateGoodsTraining(l.ctx, l.svcCtx, tsGoods, in.UserId)
			if err != nil {
				l.Logger.Error("修改商品状态失败", tsGoods)
				return nil, err
			}
			//将筛选出的商品添加到训练商品列表
			trainingGoodsList = append(trainingGoodsList, tsGoods)
		}
	}

	// 请求获取商品JSON
	err = ApplyGoodsJson(l.svcCtx, trainingGoodsList)
	if err != nil {
		l.Logger.Info("发送获取商品JSON请求失败", err)
		return nil, err
	}

	// 等待2分钟
	time.Sleep(time.Minute * 2)

	// 每6分钟调用一次接口 连续10次失败则结束
	FetchAndSaveGoodsJson(l.Logger, l.ctx, l.svcCtx, trainingGoodsList)

	// 最终保存到ES的结果文档
	var goodsDocumentList []*PddGoodsDocument
	// 获取并解析商品JSON到结果文档列表
	GetAndParseGoodsJson(l.Logger, tsShop, goodsDocumentList, trainingGoodsList)

	// 发起店铺训练批处理 获取返回的batchId
	var createBatchTaskResp string
	batchId, err := CreateBatchTask(l.Logger, l.svcCtx, goodsDocumentList, &createBatchTaskResp)
	if err != nil {
		l.Logger.Error("发送创建店铺训练批处理请求失败", err)
		return nil, err
	}

	// 等待2分钟
	time.Sleep(time.Minute * 2)

	// 轮询等待批处理完成 获取返回的fileId
	fileId, err := GetBatchTaskStatus(l.Logger, l.svcCtx, batchId)
	if err != nil {
		return nil, err
	}

	// 获取批处理结果 对于识别失败的结果将不返回
	var batchTaskResultResp string
	err = utils.HTTPGetAndParseJSON(l.svcCtx.Config.TrainingGoodsConf.QueryBatchTaskResultUrl+"?file_id="+fileId, &batchTaskResultResp)
	if err != nil {
		l.Logger.Error("发送获取批处理结果请求失败", err)
		return nil, err
	}

	// 解析结果写入goodsDocument
	var batchTaskResultMap map[string]*gjson.Result
	for _, batchTaskResult := range gjson.Get(batchTaskResultResp, "data").Array() {
		batchTaskResultMap[batchTaskResult.Get("custom_id").String()] = &batchTaskResult
	}
	for _, goodsDocument := range goodsDocumentList {
		// 只有训练成功的商品才去获取训练结果
		if batchTaskResultMap[goodsDocument.PlatformGoodsId] != nil {
			// 保存训练结果和消耗的token
			goodsDocument.DetailGalleryDescription = batchTaskResultMap[goodsDocument.PlatformGoodsId].Get("content").String()
			goodsDocument.Token = batchTaskResultMap[goodsDocument.PlatformGoodsId].Get("token").Int()
		}
		// 保存训练结果到ES
		es := l.svcCtx.Elasticsearch
		res, err := es.Index().Index("training_goods").BodyJson(goodsDocument).Refresh("true").Do(context.Background())
		if err != nil {
			logx.Errorf("商品解析结果写入ES失败, err :%s", err.Error())
			continue
		}
		logx.Infof("商品解析结果写入ES成功, res :%v", res)
	}

	// 更新数据库状态为训练完成 同时保存训练历史
	err = UpdateShopTrainingComplete(l.ctx, l.svcCtx, tsShop, in.UserId)
	if err != nil {
		l.Logger.Error("修改店铺状态失败", in)
		return nil, err
	}
	for _, trainingGoods := range trainingGoodsList {
		err = UpdateGoodsTrainingComplete(l.ctx, l.svcCtx, trainingGoods, in.UserId)
		if err != nil {
			l.Logger.Error("修改商品状态失败", trainingGoods)
			return nil, err
		}
	}

	// 返回正常
	return &training.TrainingShopResp{}, nil
}

func ApplyGoodsJson(svcCtx *svc.ServiceContext, trainingGoodsList []*orm.TsGoods) error {
	var skuList []*string
	for _, trainingGoods := range trainingGoodsList {
		skuList = append(skuList, &trainingGoods.PlatformId)
	}
	// 发送获取商品JSON请求
	var applyGoodsJsonResp string
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

func GetAndParseGoodsJson(log logx.Logger, tsShop *orm.TsShop, goodsDocumentList []*PddGoodsDocument, trainingGoodsList []*orm.TsGoods) {
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

func CreateBatchTask(log logx.Logger, svcCtx *svc.ServiceContext, goodsDocumentList []*PddGoodsDocument, createBatchTaskResp *string) (string, error) {
	var batchImages []*ImageInfo
	for _, goodsDocument := range goodsDocumentList {
		batchImages = append(batchImages, &ImageInfo{
			ID:   goodsDocument.PlatformGoodsId,
			URLs: goodsDocument.PictureUrlList,
		})
	}
	err := utils.HTTPPostAndParseJSON(svcCtx.Config.TrainingGoodsConf.CreateBatchTaskUrl, CreateBatchTaskReq{
		SystemPrompt: "what do you see ？ reply in Chinese",
		BatchImages:  batchImages,
	}, createBatchTaskResp)
	if err != nil {
		log.Error("发送创建店铺训练批处理请求失败", err)
		return "", err
	}
	batchId := gjson.Get(*createBatchTaskResp, "data.response.batch_info.id")
	return batchId.String(), nil
}

func GetBatchTaskStatus(log logx.Logger, svcCtx *svc.ServiceContext, batchId string) (string, error) {
	var fileId string
	for {
		var batchTaskStatusResp string
		err := utils.HTTPGetAndParseJSON(svcCtx.Config.TrainingGoodsConf.QueryBatchTaskStatusUrl+"?batch_id="+batchId, &batchTaskStatusResp)
		if err != nil {
			log.Error("发送获取批处理状态请求失败", err)
			return "nil", err
		}
		status := gjson.Get(batchTaskStatusResp, "data.status")
		if status.String() == "completed" {
			fileId = gjson.Get(batchTaskStatusResp, "data.output_file_id").String()
			break
		}
		time.Sleep(time.Minute * 2)
	}
	return fileId, nil
}

func UpdateShopTraining(ctx context.Context, svcCtx *svc.ServiceContext, tsShop *orm.TsShop, userId int64) error {
	tsShop.TrainingStatus = 1
	tsShop.TrainingTimes += 1
	tsShop.UpdateTime = time.Now()
	tsShop.UpdateBy = userId
	err := svcCtx.TsShopModel.Update(ctx, tsShop)
	if err != nil {
		return err
	}
	return nil
}

func UpdateShopTrainingComplete(ctx context.Context, svcCtx *svc.ServiceContext, tsShop *orm.TsShop, userId int64) error {
	tsShop.TrainingStatus = 2
	tsShop.UpdateTime = time.Now()
	tsShop.UpdateBy = userId
	err := svcCtx.TsShopModel.Update(ctx, tsShop)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGoodsTraining(ctx context.Context, svcCtx *svc.ServiceContext, tsGoods *orm.TsGoods, userId int64) error {
	tsGoods.TrainingStatus = 1
	tsGoods.TrainingTimes += 1
	tsGoods.UpdateTime = time.Now()
	tsGoods.UpdateBy = userId
	tsGoods.GoodsJsonUrl = "" //每次训练开始 把获取商品json的url字段置空
	err := svcCtx.TsGoodsModel.Update(ctx, tsGoods)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGoodsTrainingComplete(ctx context.Context, svcCtx *svc.ServiceContext, tsGoods *orm.TsGoods, userId int64) error {
	tsGoods.TrainingStatus = 2
	tsGoods.UpdateTime = time.Now()
	tsGoods.UpdateBy = userId
	err := svcCtx.TsGoodsModel.Update(ctx, tsGoods)
	if err != nil {
		return err
	}
	return nil
}

// ParsePddGoods 解析拼多多商品JSON
func ParsePddGoods(goodsJson string, tsShop *orm.TsShop, tsGoods *orm.TsGoods) *PddGoodsDocument {
	// 店铺id
	mallId := gjson.Get(goodsJson, "mall_entrance.mall_data.mall_id")
	// 商品sku标签列表
	var skuSpecsMap map[string]string
	// 商品中的图片列表
	var pictureUrlList []string
	for _, sku := range gjson.Get(goodsJson, "sku").Array() {
		pictureUrlList = append(pictureUrlList, sku.Get("thumb_url").String())
		for _, skuSpec := range sku.Get("specs").Array() {
			skuSpecsMap[skuSpec.Get("spec_key").String()] = skuSpec.Get("spec_value").String()
		}
	}
	for _, gallery := range gjson.Get(goodsJson, "goods.gallery").Array() {
		pictureUrlList = append(pictureUrlList, gallery.Get("url").String())
	}
	// 商品服务承诺列表
	var ServicePromiseMap map[string]string
	for _, servicePromise := range gjson.Get(goodsJson, "service_promise").Array() {
		ServicePromiseMap[servicePromise.Get("type").String()] = servicePromise.Get("desc").String()
	}
	// 团购价格和基础价格
	groupPrice, _ := strconv.ParseFloat(gjson.Get(goodsJson, "price.min_group_price").String(), 64)
	normalPrice, _ := strconv.ParseFloat(gjson.Get(goodsJson, "price.max_normal_price").String(), 64)
	// 卖点
	sellPointTagList := gjson.Get(goodsJson, "ui.carousel_section.sell_point_tag_list").Array()
	// 商品提示
	promptExplain := gjson.Get(goodsJson, "goods.prompt_explain").String()
	goodsDocument := &PddGoodsDocument{
		ShopId:  tsShop.Id,
		GoodsId: tsGoods.Id,
		Uuid:    tsShop.Uuid,
		UserId:  tsShop.UserId,

		PlatformMallId:  mallId.String(),
		PlatformGoodsId: tsGoods.PlatformId,
		GoodsUrl:        tsGoods.GoodsUrl,
		GoodsJson:       goodsJson,
		GoodsName:       tsGoods.GoodsName,

		SkuSpecs:                 skuSpecsMap,
		GroupPrice:               groupPrice,
		NormalPrice:              normalPrice,
		ServicePromise:           ServicePromiseMap,
		SellPointTagList:         sellPointTagList,
		PromptExplain:            promptExplain,
		DetailGalleryDescription: "",

		PictureUrlList: pictureUrlList,
		Token:          0,
		CreatedAt:      time.Now(),
	}
	return goodsDocument
}
