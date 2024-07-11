package model

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ ShoppresettingshoptitlesModel = (*customShoppresettingshoptitlesModel)(nil)

type (
	// ShoppresettingshoptitlesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShoppresettingshoptitlesModel.
	ShoppresettingshoptitlesModel interface {
		shoppresettingshoptitlesModel
		FindNewOneByUuidAndUserId(ctx context.Context, uuid string, userId int64) (*Shoppresettingshoptitles, error)
	}

	customShoppresettingshoptitlesModel struct {
		*defaultShoppresettingshoptitlesModel
	}
)

// NewShoppresettingshoptitlesModel returns a model for the mongo.
func NewShoppresettingshoptitlesModel(url, db, collection string) ShoppresettingshoptitlesModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customShoppresettingshoptitlesModel{
		defaultShoppresettingshoptitlesModel: newDefaultShoppresettingshoptitlesModel(conn),
	}
}

func (m *defaultShoppresettingshoptitlesModel) FindNewOneByUuidAndUserId(ctx context.Context, uuid string, userId int64) (*Shoppresettingshoptitles, error) {
	var data Shoppresettingshoptitles
	opts := options.FindOne().SetSort(map[string]interface{}{"createdAt": -1})
	err := m.conn.FindOne(ctx, &data, bson.M{"uuid": uuid, "userId": userId}, opts)
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, mon.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
