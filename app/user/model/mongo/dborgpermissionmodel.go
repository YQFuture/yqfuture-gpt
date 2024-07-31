package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ DborgpermissionModel = (*customDborgpermissionModel)(nil)

type (
	// DborgpermissionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDborgpermissionModel.
	DborgpermissionModel interface {
		dborgpermissionModel
		InsertOne(ctx context.Context, data *Dborgpermission) (*mongo.InsertOneResult, error)
	}

	customDborgpermissionModel struct {
		*defaultDborgpermissionModel
	}
)

// NewDborgpermissionModel returns a model for the mongo.
func NewDborgpermissionModel(url, db, collection string) DborgpermissionModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customDborgpermissionModel{
		defaultDborgpermissionModel: newDefaultDborgpermissionModel(conn),
	}
}

func (m *customDborgpermissionModel) InsertOne(ctx context.Context, data *Dborgpermission) (*mongo.InsertOneResult, error) {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	result, err := m.conn.InsertOne(ctx, data)
	return result, err
}
