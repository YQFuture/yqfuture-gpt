package model

import (
	"time"
	"yufuture-gpt/app/training/model/common"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shoppresettingshoptitles struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UpdateAt time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`

	ShopId     int64  `bson:"shopId" json:"shopId"`
	PlatformId string `bson:"platformId" json:"platformId"`
	UUID       string `bson:"uuid" json:"uuid"`
	UserID     int64  `bson:"userId" json:"userId"`

	PreSettingToken    int64                      `bson:"preSettingToken" json:"preSettingToken"`
	PresettingHashRate int64                      `bson:"presettingHashRate" json:"presettingHashRate"`
	PreSettingTime     time.Time                  `bson:"preSettingTime" json:"preSettingTime"`
	GoodsDocument      []*common.PddGoodsDocument `bson:"goodsDocument" json:"goodsDocument"`
}
