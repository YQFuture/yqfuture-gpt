package common

import "time"

// PddGoodsDocument 拼多多商品训练结果
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
