package model

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ DbuseroperationModel = (*customDbuseroperationModel)(nil)

type (
	// DbuseroperationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDbuseroperationModel.
	DbuseroperationModel interface {
		dbuseroperationModel
		FindPageListByUserIdAndOrgId(ctx context.Context, userId, orgId, pageNum, pageSize int64, queryString string) (*[]*Dbuseroperation, error)
	}

	customDbuseroperationModel struct {
		*defaultDbuseroperationModel
	}
)

// NewDbuseroperationModel returns a model for the mongo.
func NewDbuseroperationModel(url, db, collection string) DbuseroperationModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customDbuseroperationModel{
		defaultDbuseroperationModel: newDefaultDbuseroperationModel(conn),
	}
}

func (m *customDbuseroperationModel) FindPageListByUserIdAndOrgId(ctx context.Context, userId, orgId, pageNum, pageSize int64, queryString string) (*[]*Dbuseroperation, error) {
	// 创建查询条件
	filter := bson.M{}
	if userId != 0 {
		filter["userId"] = userId
	}
	if orgId != 0 {
		filter["orgId"] = orgId
	}
	if queryString != "" {
		// 添加模糊查询条件
		filter["operationdesc"] = bson.M{"$regex": queryString, "$options": "i"}
	}

	// 设置排序选项
	sortOptions := options.Find().SetSort(bson.D{{"createat", -1}})

	// 设置分页选项
	skip := (pageNum - 1) * pageSize
	limitOption := options.Find().SetLimit(pageSize).SetSkip(skip)

	// 执行查询
	var data []*Dbuseroperation
	err := m.conn.Find(ctx, &data, filter, limitOption.SetSort(sortOptions))
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, mon.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}
