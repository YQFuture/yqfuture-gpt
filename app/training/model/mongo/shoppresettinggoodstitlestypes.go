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

	GoodsId    int64  `bson:"goodsId" json:"goodsId"`
	PlatformId string `bson:"platformId" json:"platformId"`
	UserID     int64  `bson:"userId" json:"userId"`

	PreSettingToken    int64                      `bson:"preSettingToken" json:"preSettingToken"`
	PresettingPower    int64                      `bson:"presettingPower" json:"presettingPower"`
	PresettingFileSize int64                      `bson:"presettingFileSize" json:"presettingFileSize"`
	PreSettingTime     time.Time                  `bson:"preSettingTime" json:"preSettingTime"`
	GoodsDocumentList  []*common.PddGoodsDocument `bson:"goodsDocumentList" json:"goodsDocumentList"`
}
