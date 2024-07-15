package model

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ DbpresettingshoptitlesModel = (*customDbpresettingshoptitlesModel)(nil)

type (
	// DbpresettingshoptitlesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDbpresettingshoptitlesModel.
	DbpresettingshoptitlesModel interface {
		dbpresettingshoptitlesModel
		FindNewOneByUuidAndUserId(ctx context.Context, uuid string, userId int64) (*Dbpresettingshoptitles, error)
	}

	customDbpresettingshoptitlesModel struct {
		*defaultDbpresettingshoptitlesModel
	}
)

// NewDbpresettingshoptitlesModel returns a model for the mongo.
func NewDbpresettingshoptitlesModel(url, db, collection string) DbpresettingshoptitlesModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customDbpresettingshoptitlesModel{
		defaultDbpresettingshoptitlesModel: newDefaultDbpresettingshoptitlesModel(conn),
	}
}

func (m *defaultDbpresettingshoptitlesModel) FindNewOneByUuidAndUserId(ctx context.Context, uuid string, userId int64) (*Dbpresettingshoptitles, error) {
	var data Dbpresettingshoptitles
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
