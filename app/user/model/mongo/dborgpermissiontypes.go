package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Dborgpermission struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	PermissionList []*Permission `bson:"permissionList,omitempty" json:"permissionList,omitempty"`
	RoleList       []*Role       `bson:"roleList,omitempty" json:"roleList,omitempty"`
	UserList       []*User       `bson:"userList,omitempty" json:"userList,omitempty"`

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}

type Permission struct {
	Id         int64  `bson:"id,omitempty" json:"id,omitempty"`
	ParentId   int64  `bson:"parentId,omitempty" json:"parentId,omitempty"`
	Name       string `bson:"name,omitempty" json:"name,omitempty"`
	Perm       string `bson:"perm,omitempty" json:"perm,omitempty"`
	Url        string `bson:"url,omitempty" json:"url,omitempty"`
	Sort       int64  `bson:"sort,omitempty" json:"sort,omitempty"`
	TemplateId int64  `bson:"templateId,omitempty" json:"templateId,omitempty"`
	ResourceId int64  `bson:"resourceId,omitempty" json:"resourceId,omitempty"`
}

type Role struct {
	Id             int64    `bson:"id,omitempty" json:"id,omitempty"`
	Name           string   `bson:"name,omitempty" json:"name,omitempty"`
	PermissionList []*int64 `bson:"permissionList,omitempty" json:"permissionList,omitempty"`
}

type User struct {
	Id                       int64    `bson:"id,omitempty" json:"id,omitempty"`
	RoleList                 []*int64 `bson:"roleList,omitempty" json:"roleList,omitempty"`
	KeywordSwitchingShopList []*int64 `bson:"keywordSwitchingShopList,omitempty" json:"keywordSwitchingShopList,omitempty"` // 关键词转接店铺列表
	ExceptionDutyShopList    []*int64 `bson:"exceptionDutyShopList,omitempty" json:"exceptionDutyShopList,omitempty"`       // 异常责任店铺列表
}
