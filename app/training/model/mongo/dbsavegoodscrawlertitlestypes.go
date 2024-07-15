package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dbsavegoodscrawlertitles struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	SerialNumber string         `bson:"serialNumber" json:"serialNumber"`
	ShopId       int64          `bson:"shopId" json:"shopId"`
	UserID       int64          `bson:"userId" json:"userId"`
	UUID         string         `bson:"uuid" json:"uuid"`
	Platform     int64          `bson:"platform" json:"platform"`
	ShopName     string         `bson:"shopName" json:"shopName"`
	GoodsList    []*GoodsIdList `bson:"goodsList" json:"goodsList"`

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}

type GoodsIdList struct {
	GoodsId    int64  `bson:"goodsId" json:"goodsId"`
	PlatformId string `bson:"platformId" json:"platformId"`
	Url        string `bson:"url" json:"url"`
}
