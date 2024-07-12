package common

import "time"

// PddGoodsDocument 拼多多商品训练结果
type PddGoodsDocument struct {
	ShopId  int64  `bson:"shopId" json:"shopId"`   // 对应mysql店铺id
	GoodsId int64  `bson:"goodsId" json:"goodsId"` // 对应mysql商品id
	Uuid    string `bson:"uuid" json:"uuid"`
	UserId  int64  `bson:"userId" json:"userId"`

	PlatformMallId  string `bson:"platformMallId" json:"platformMallId"`   // 对应json店铺id
	PlatformGoodsId string `bson:"platformGoodsId" json:"platformGoodsId"` // 对应json商品id
	GoodsUrl        string `bson:"goodsUrl" json:"goodsUrl"`
	GoodsJson       string `bson:"goodsJson" json:"goodsJson"`
	GoodsName       string `bson:"goodsName" json:"goodsName"`

	SkuSpecs                 map[string]string `bson:"skuSpecs" json:"skuSpecs"`
	GroupPrice               float64           `bson:"groupPrice" json:"groupPrice"`
	NormalPrice              float64           `bson:"normalPrice" json:"normalPrice"`
	ServicePromise           map[string]string `bson:"servicePromise" json:"servicePromise"`                     // 商品服务承诺列表
	SellPointTagList         interface{}       `bson:"sellPointTagList" json:"sellPointTagList"`                 // 卖点
	PromptExplain            string            `bson:"promptExplain" json:"promptExplain"`                       // 商品提示
	DetailGalleryDescription string            `bson:"detailGalleryDescription" json:"detailGalleryDescription"` // 图片训练结果描述

	PictureUrlList []string  `bson:"pictureUrlList" json:"pictureUrlList"`
	Token          int64     `bson:"token" json:"token"`       // 消耗的token
	Power          int64     `bson:"power" json:"power"`       // 消耗的算力
	FileSize       int64     `bson:"fileSize" json:"fileSize"` // 文件大小
	CreatedAt      time.Time `bson:"createdAt" json:"createdAt"`
}
