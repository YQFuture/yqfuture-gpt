package model

import (
	"time"
	"yufuture-gpt/app/training/model/common"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shoppresettinggoodstitles struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UpdateAt time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`

	GoodsId    string `bson:"goodsId" json:"goodsId"`
	PlatformId string `bson:"platformId" json:"platformId"`

	PreSettingToken    int64                      `bson:"preSettingToken" json:"preSettingToken"`
	PresettingHashRate int64                      `bson:"presettingHashRate" json:"presettingHashRate"`
	PreSettingTime     time.Time                  `bson:"preSettingTime" json:"preSettingTime"`
	GoodsDocumentList  []*common.PddGoodsDocument `bson:"goodsDocumentList" json:"goodsDocumentList"`
}
