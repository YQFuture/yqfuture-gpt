package thirdparty

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/model/common"
	"yufuture-gpt/app/training/model/orm"
	"yufuture-gpt/common/consts"
	"yufuture-gpt/common/utils"
)

type ApplyGoodsIdListReq struct {
	SerialNumber string        `json:"serialNumber"`
	Uuid         string        `json:"uuid"`
	ShopName     string        `json:"shop_name"`
	Platform     string        `json:"platform"`
	Token        string        `json:"token"`
	Cookies      []interface{} `json:"cookies"`
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
