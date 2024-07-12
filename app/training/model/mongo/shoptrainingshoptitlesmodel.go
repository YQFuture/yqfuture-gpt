package model

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ ShoptrainingshoptitlesModel = (*customShoptrainingshoptitlesModel)(nil)

type (
	// ShoptrainingshoptitlesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShoptrainingshoptitlesModel.
	ShoptrainingshoptitlesModel interface {
		shoptrainingshoptitlesModel
		FindNewOneByUuidAndUserId(ctx context.Context, uuid string, userId int64) (*Shoptrainingshoptitles, error)
		FindOneBySerialNumber(ctx context.Context, serialNumber string) (*Shoptrainingshoptitles, error)
	}

	customShoptrainingshoptitlesModel struct {
		*defaultShoptrainingshoptitlesModel
	}
)

// NewShoptrainingshoptitlesModel returns a model for the mongo.
func NewShoptrainingshoptitlesModel(url, db, collection string) ShoptrainingshoptitlesModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customShoptrainingshoptitlesModel{
		defaultShoptrainingshoptitlesModel: newDefaultShoptrainingshoptitlesModel(conn),
	}
}

func (m *defaultShoptrainingshoptitlesModel) FindNewOneByUuidAndUserId(ctx context.Context, uuid string, userId int64) (*Shoptrainingshoptitles, error) {
	var data Shoptrainingshoptitles
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

func (m *defaultShoptrainingshoptitlesModel) FindOneBySerialNumber(ctx context.Context, serialNumber string) (*Shoptrainingshoptitles, error) {
	var data Shoptrainingshoptitles
	err := m.conn.FindOne(ctx, &data, bson.M{"serialNumber": serialNumber})
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, mon.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
