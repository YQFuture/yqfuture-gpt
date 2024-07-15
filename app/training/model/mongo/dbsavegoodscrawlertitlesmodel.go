package model

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ DbsavegoodscrawlertitlesModel = (*customDbsavegoodscrawlertitlesModel)(nil)

type (
	// DbsavegoodscrawlertitlesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDbsavegoodscrawlertitlesModel.
	DbsavegoodscrawlertitlesModel interface {
		dbsavegoodscrawlertitlesModel
		FindNewOneByUuidAndUserId(ctx context.Context, uuid string, userId int64) (*Dbsavegoodscrawlertitles, error)
		FindOneBySerialNumber(ctx context.Context, serialNumber string) (*Dbsavegoodscrawlertitles, error)
	}

	customDbsavegoodscrawlertitlesModel struct {
		*defaultDbsavegoodscrawlertitlesModel
	}
)

// NewDbsavegoodscrawlertitlesModel returns a model for the mongo.
func NewDbsavegoodscrawlertitlesModel(url, db, collection string) DbsavegoodscrawlertitlesModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customDbsavegoodscrawlertitlesModel{
		defaultDbsavegoodscrawlertitlesModel: newDefaultDbsavegoodscrawlertitlesModel(conn),
	}
}

func (m *defaultDbsavegoodscrawlertitlesModel) FindNewOneByUuidAndUserId(ctx context.Context, uuid string, userId int64) (*Dbsavegoodscrawlertitles, error) {
	var data Dbsavegoodscrawlertitles
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

func (m *defaultDbsavegoodscrawlertitlesModel) FindOneBySerialNumber(ctx context.Context, serialNumber string) (*Dbsavegoodscrawlertitles, error) {
	var data Dbsavegoodscrawlertitles
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
