package thirdparty

import (
	"github.com/tidwall/gjson"
	"strconv"
	"time"
	"yufuture-gpt/app/training/model/common"
	"yufuture-gpt/app/training/model/orm"
)

func ParsePddGoods(goodsJson string, tsShop *orm.TsShop, tsGoods *orm.TsGoods) *common.PddGoodsDocument {
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
	goodsDocument := &common.PddGoodsDocument{
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
		UpdatedAt:      time.Now(),
	}
	return goodsDocument
}
