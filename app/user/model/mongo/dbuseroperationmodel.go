package model

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ DbuseroperationModel = (*customDbuseroperationModel)(nil)

type (
	// DbuseroperationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDbuseroperationModel.
	DbuseroperationModel interface {
		dbuseroperationModel
		FindPageListByUserIdAndOrgId(ctx context.Context, userId, orgId, pageNum, pageSize, startTime, endTime int64, queryString string) (*[]*Dbuseroperation, error)
		FindPageTotalByUserIdAndOrgId(ctx context.Context, userId, orgId, startTime, endTime int64, queryString string) (int64, error)
		FindListByUserIdAndOrgId(ctx context.Context, userId, orgId, startTime, endTime int64, queryString string) (*[]*Dbuseroperation, error)
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

func (m *customDbuseroperationModel) FindPageListByUserIdAndOrgId(ctx context.Context, userId, orgId, pageNum, pageSize, startTime, endTime int64, queryString string) (*[]*Dbuseroperation, error) {
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

	// 将开始和结束时间从 int64 转换为 time.Time 类型
	startTimeTime := time.UnixMilli(startTime)
	endTimeTime := time.UnixMilli(endTime)

	// 添加时间范围条件
	if startTime != 0 && endTime != 0 {
		filter["createat"] = bson.M{"$gte": startTimeTime, "$lte": endTimeTime}
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

func (m *customDbuseroperationModel) FindPageTotalByUserIdAndOrgId(ctx context.Context, userId, orgId, startTime, endTime int64, queryString string) (int64, error) {
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

	// 将开始和结束时间从 int64 转换为 time.Time 类型
	startTimeTime := time.UnixMilli(startTime)
	endTimeTime := time.UnixMilli(endTime)

	// 添加时间范围条件
	if startTime != 0 && endTime != 0 {
		filter["createat"] = bson.M{"$gte": startTimeTime, "$lte": endTimeTime}
	}

	// 执行计数操作
	count, err := m.conn.CountDocuments(ctx, filter)
	switch {
	case err == nil:
		return count, nil
	case errors.Is(err, mon.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}

func (m *customDbuseroperationModel) FindListByUserIdAndOrgId(ctx context.Context, userId, orgId, startTime, endTime int64, queryString string) (*[]*Dbuseroperation, error) {
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

	// 将开始和结束时间从 int64 转换为 time.Time 类型
	startTimeTime := time.UnixMilli(startTime)
	endTimeTime := time.UnixMilli(endTime)

	// 添加时间范围条件
	if startTime != 0 && endTime != 0 {
		filter["createat"] = bson.M{"$gte": startTimeTime, "$lte": endTimeTime}
	}

	// 设置排序选项
	sortOptions := options.Find().SetSort(bson.D{{"createat", -1}})

	// 执行查询
	var data []*Dbuseroperation
	err := m.conn.Find(ctx, &data, filter, sortOptions)
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, mon.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}
