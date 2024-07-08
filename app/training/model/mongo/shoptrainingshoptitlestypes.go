package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shoptrainingshoptitles struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UpdateAt time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`

	UserID    int32  `bson:"userId" json:"userId"`
	UUID      string `bson:"uuid" json:"uuid"`
	Platform  int64  `bson:"platform" json:"platform"`
	ShopName  string `bson:"shop_name" json:"shopName"`
	GoodsList []struct {
		PlatformId string `bson:"id" json:"platformId"`
		Url        string `bson:"url" json:"url"`
	} `bson:"goods_list" json:"goodsList"`
}
