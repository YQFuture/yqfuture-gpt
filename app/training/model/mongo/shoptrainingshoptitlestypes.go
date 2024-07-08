package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shoptrainingshoptitles struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UpdateAt time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`

	ShopId    int64                    `bson:"shopId" json:"shopId"`
	UserID    int64                    `bson:"userId" json:"userId"`
	UUID      string                   `bson:"uuid" json:"uuid"`
	Platform  int64                    `bson:"platform" json:"platform"`
	ShopName  string                   `bson:"shop_name" json:"shopName"`
	GoodsList []*ShopTrainingGoodsList `bson:"goods_list" json:"goodsList"`
}

type ShopTrainingGoodsList struct {
	GoodsId    int64  `bson:"id" json:"goodsId"`
	PlatformId string `bson:"id" json:"platformId"`
	Url        string `bson:"url" json:"url"`
}
