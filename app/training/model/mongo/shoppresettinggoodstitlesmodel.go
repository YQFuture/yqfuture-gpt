package model

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ ShoppresettinggoodstitlesModel = (*customShoppresettinggoodstitlesModel)(nil)

type (
	// ShoppresettinggoodstitlesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShoppresettinggoodstitlesModel.
	ShoppresettinggoodstitlesModel interface {
		shoppresettinggoodstitlesModel
		FindNewOneByGoodsId(ctx context.Context, goodsId int64) (*Shoppresettinggoodstitles, error)
	}

	customShoppresettinggoodstitlesModel struct {
		*defaultShoppresettinggoodstitlesModel
	}
)

// NewShoppresettinggoodstitlesModel returns a model for the mongo.
func NewShoppresettinggoodstitlesModel(url, db, collection string) ShoppresettinggoodstitlesModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customShoppresettinggoodstitlesModel{
		defaultShoppresettinggoodstitlesModel: newDefaultShoppresettinggoodstitlesModel(conn),
	}
}

func (m *defaultShoppresettinggoodstitlesModel) FindNewOneByGoodsId(ctx context.Context, goodsId int64) (*Shoppresettinggoodstitles, error) {
	var data Shoppresettinggoodstitles
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
