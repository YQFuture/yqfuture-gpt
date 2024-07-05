package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type YqfutureShop struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    int32              `bson:"userId"`
	UUID      string             `bson:"uuid"`
	Platform  int64              `bson:"platform"`
	ShopName  string             `bson:"shop_name"`
	GoodsList []struct {
		PlatformId string `bson:"platform_id"`
		Url        string `bson:"url"`
	} `bson:"goods"`
	CreateTime time.Time
	UpdateTime time.Time
}
