package model

import (
	"time"
	"yufuture-gpt/app/training/model/common"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dbpresettingshoptitles struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	ShopId             int64                      `bson:"shopId" json:"shopId"`                         // 店铺id
	PlatformId         string                     `bson:"platformId" json:"platformId"`                 // 平台店铺id
	UUID               string                     `bson:"uuid" json:"uuid"`                             // 店铺uuid
	UserID             int64                      `bson:"userId" json:"userId"`                         // 用户id
	PreSettingToken    int64                      `bson:"preSettingToken" json:"preSettingToken"`       // 预设token
	PresettingPower    int64                      `bson:"presettingPower" json:"presettingPower"`       // 预设算力
	PresettingFileSize int64                      `bson:"presettingFileSize" json:"presettingFileSize"` // 预设文件大小
	PreSettingTime     time.Time                  `bson:"preSettingTime" json:"preSettingTime"`         // 预设时间
	GoodsDocumentList  []*common.PddGoodsDocument `bson:"goodsDocumentList" json:"goodsDocumentList"`   // 商品文档列表

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
