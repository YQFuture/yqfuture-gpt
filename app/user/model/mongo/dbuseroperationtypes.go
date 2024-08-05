package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dbuseroperation struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId int64              `bson:"userId,omitempty" json:"userId,omitempty"`
	OrgId  int64              `bson:"orgId,omitempty" json:"orgId,omitempty"`

	ResourceType string `bson:"resourceType,omitempty" json:"resourceType,omitempty"` // 资源类型 0 未定义 1 店铺
	ResourceId   int64  `bson:"resourceId,omitempty" json:"resourceId,omitempty"`

	OperationDesc string `bson:"operationDesc,omitempty" json:"operationDesc,omitempty"` // 操作描述

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
