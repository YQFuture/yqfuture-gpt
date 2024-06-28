package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Kfgptaccountsentities struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UpdateAt   time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt   time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
	UserID     int32              `bson:"userId"`
	UUID       string             `bson:"uuid"`
	Type       int32              `bson:"type"`
	Account    string             `bson:"account"`
	Password   string             `bson:"password"`
	NickName   string             `bson:"nick_name"`
	ShopName   string             `bson:"shop_name"`
	Avatar     string             `bson:"avatar"`
	Remark     string             `bson:"remark"`
	GroupNames []string           `bson:"group_names"`
	Status     int32              `bson:"status"`
	Token      string             `bson:"token"`
	TopTime    int32              `bson:"top_time"`
	CreateTime time.Time
	UpdateTime time.Time
	Version    int32 `bson:"__v"`
}
