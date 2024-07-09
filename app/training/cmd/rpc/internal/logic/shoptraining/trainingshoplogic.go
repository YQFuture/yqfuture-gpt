package shoptraininglogic

import (
	"context"
	"github.com/tidwall/gjson"
	"github.com/zeromicro/go-zero/core/logx"
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
	l.Logger.Error("", saveShop)

	// 根据uuid和userid查找出店铺
	tsShop, err := l.svcCtx.TsShopModel.FindOneByUuidAndUserId(l.ctx, in.UserId, in.Uuid)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找店铺失败", err)
		return nil, err
	}

	// 根据店铺shopId查找出enabled字段为2启用商品列表
	tsGoodsList, err := l.svcCtx.TsGoodsModel.FindEnabledListByShopId(l.ctx, tsShop.Id)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找商品失败", err)
		return nil, err
	}
	var tsGoodsMap map[string]*orm.TsGoods
	if tsGoodsList != nil {
		tsGoodsMap = make(map[string]*orm.TsGoods)
		for _, tsGoods := range *tsGoodsList {
			tsGoodsMap[tsGoods.PlatformId] = tsGoods
		}
	}

	// 更新店铺状态为训练中 添加训练次数
	tsShop.TrainingStatus = 1
	tsShop.TrainingTimes += 1
	tsShop.UpdateTime = time.Now()
	tsShop.UpdateBy = in.UserId
	err = l.svcCtx.TsShopModel.Update(l.ctx, tsShop)
	if err != nil {
		l.Logger.Error("修改店铺状态失败", in)
		return nil, err
	}

	// trainingGoodsList 本次需要训练的商品列表
	var trainingGoodsList []*orm.TsGoods
	for _, saveGoods := range saveShop.GoodsList {
		if tsGoods, ok := tsGoodsMap[saveGoods.PlatformId]; ok {
			// 同时在mongo中并且enabled字段为2启用的商品即为本次需要训练的商品
			// 更新数据库的状态
			tsGoods.TrainingStatus = 1
			tsGoods.TrainingTimes += 1
			tsGoods.UpdateTime = time.Now()
			tsGoods.UpdateBy = in.UserId
			tsGoods.GoodsJsonUrl = "" //每次训练开始 把获取商品json的url字段置空
			err = l.svcCtx.TsGoodsModel.Update(l.ctx, tsGoods)
			if err != nil {
				l.Logger.Error("修改商品状态失败", tsGoods)
				return nil, err
			}
			//将筛选出的商品添加到训练商品列表
			trainingGoodsList = append(trainingGoodsList, tsGoods)
		}
	}

	// 构建获取商品JSON请求体
	var skuList []*string
	for _, trainingGoods := range trainingGoodsList {
		skuList := append(skuList, &trainingGoods.PlatformId)
		l.Logger.Info("获取JSON的数据", skuList)
	}

	// 发送获取商品JSON请求

	// 等待2分钟

	// 每6分钟调用一次接口 连续10次失败则结束
	i := 0
	for i < 10 {
		// 调用接口

		// 将返回的获取商品json的url更新进mysql

		// 查询mysql判断是否完成 同时将完成的新获取json的url保存进training_url
		complete := true
		for _, trainingGoods := range trainingGoodsList {
			// 已经更新获取json的url的商品跳过循环
			if trainingGoods.GoodsUrl != "" {
				continue
			}
			mysqlGoods, err := l.svcCtx.TsGoodsModel.FindOne(l.ctx, trainingGoods.Id)
			if err != nil {
				l.Logger.Error("查询训练中的商品数据失败", err)
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

	// 解析商品JSON
	for _, trainingGoods := range trainingGoodsList {
		//根据获取商品JSON的url获取商品JSON
		if trainingGoods.GoodsJsonUrl == "" {
			continue
		}
		var goodsJson string
		err := utils.HTTPGetAndParseJSON(trainingGoods.GoodsUrl, &goodsJson)
		if err != nil {
			l.Logger.Error("根据url获取商品json数据失败", err)
			continue
		}
		if goodsJson == "" {
			l.Logger.Error("根据url获取商品json数据失败, 返回的json为空", err)
			continue
		}

		goodsDocument := parsePddGoods(l, goodsJson, tsShop, *trainingGoods)

		es := l.svcCtx.Elasticsearch
		res, err := es.Index().Index("training_goods").BodyJson(goodsDocument).Refresh("true").Do(context.Background())
		if err != nil {
			logx.Errorf("商品解析结果写入ES失败, err :%s", err.Error())
			continue
		}
		logx.Infof("商品解析结果写入ES成功, res :%v", res)

	}

	// 训练商品

	// 保存训练结果

	// 等待所有商品训练完成

	// 更新数据库状态

	// 返回正常
	return &training.TrainingShopResp{}, nil
}

type GoodsJsonResponse struct {
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

func parsePddGoods(l *TrainingShopLogic, goodsJson string, tsShop *orm.TsShop, tsGoods orm.TsGoods) *GoodsDocument {
	mallId := gjson.Get(goodsJson, "mall_entrance.mall_data.mall_id")
	skuList := gjson.Get(goodsJson, "sku")
	// 商品中的图片列表
	var pictureUrlList []string
	for _, sku := range gjson.Get(goodsJson, "sku").Array() {
		pictureUrlList = append(pictureUrlList, sku.Get("thumb_url").String())
	}
	for _, sku := range gjson.Get(goodsJson, "goods.gallery").Array() {
		pictureUrlList = append(pictureUrlList, sku.Get("url").String())
	}
	goodsDocument := &GoodsDocument{
		ShopId:  tsShop.Id,
		GoodsId: tsGoods.Id,
		Uuid:    tsShop.Uuid,
		UserId:  tsShop.UserId,

		PlatformMallId:  mallId.String(),
		PlatformGoodsId: tsGoods.PlatformId,
		GoodsUrl:        tsGoods.GoodsUrl,
		GoodsJson:       goodsJson,
		GoodsName:       tsGoods.GoodsName,

		GoodsSkus:                skuList.Array(),
		GroupPrice:               0,
		NormalPrice:              0,
		Service:                  "",
		SellPointTagList:         []string{},
		PromptExplain:            "",
		DetailGalleryDescription: "",

		PictureUrlList: pictureUrlList,
		TrainingResult: TrainingResult{},
		Token:          0,
		CreatedAt:      time.Now(),
	}

	return goodsDocument
}

// 保存到es的商品训练结果
type GoodsDocument struct {
	ShopId  int64  `json:"shop_id"`    // 对应mysql店铺id
	GoodsId int64  `json:"GoodsId_id"` // 对应mysql商品id
	Uuid    string `json:"uuid"`
	UserId  int64  `json:"user_id"`

	PlatformMallId  string `json:"platform_mall_id"`  // 对应json店铺id
	PlatformGoodsId string `json:"platform_goods_id"` // 对应json商品id
	GoodsUrl        string `json:"goods_url"`
	GoodsJson       string `json:"goods_json"`
	GoodsName       string `json:"goods_name"`

	GoodsSkus                interface{} `json:"goods_skus"`
	GroupPrice               float64     `json:"group_price"`
	NormalPrice              float64     `json:"normal_price"`
	Service                  string      `json:"service"`
	SellPointTagList         []string    `json:"sell_point_tag_list"`        // 卖点
	PromptExplain            string      `json:"prompt_explain"`             // 提示
	DetailGalleryDescription string      `json:"detail_gallery_description"` // 图片训练结果描述

	PictureUrlList []string       `json:"picture_url_list"`
	TrainingResult TrainingResult `json:"training_result"`
	Token          int            `json:"token"` // 消耗的token
	CreatedAt      time.Time      `json:"create_time"`
}

// 调用gpt的训练结果
type TrainingResult struct {
	Status bool `json:"status"`
	Data   struct {
		Response string `json:"response"` // 训练结果
		Token    int    `json:"token"`    // 消耗的token
	} `json:"data"`
	Msg string `json:"msg"`
}
