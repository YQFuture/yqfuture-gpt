package model

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ DbpresettinggoodstitlesModel = (*customDbpresettinggoodstitlesModel)(nil)

type (
	// DbpresettinggoodstitlesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDbpresettinggoodstitlesModel.
	DbpresettinggoodstitlesModel interface {
		dbpresettinggoodstitlesModel
		FindNewOneByGoodsId(ctx context.Context, goodsId int64) (*Dbpresettinggoodstitles, error)
	}

	customDbpresettinggoodstitlesModel struct {
		*defaultDbpresettinggoodstitlesModel
	}
)

// NewDbpresettinggoodstitlesModel returns a model for the mongo.
func NewDbpresettinggoodstitlesModel(url, db, collection string) DbpresettinggoodstitlesModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customDbpresettinggoodstitlesModel{
		defaultDbpresettinggoodstitlesModel: newDefaultDbpresettinggoodstitlesModel(conn),
	}
}

func (m *defaultDbpresettinggoodstitlesModel) FindNewOneByGoodsId(ctx context.Context, goodsId int64) (*Dbpresettinggoodstitles, error) {
	var data Dbpresettinggoodstitles
	opts := options.FindOne().SetSort(map[string]interface{}{"createdAt": -1})
	err := m.conn.FindOne(ctx, &data, bson.M{"goodsId": goodsId}, opts)
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, mon.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
