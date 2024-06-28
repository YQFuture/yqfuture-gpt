package model

import "github.com/zeromicro/go-zero/core/stores/mon"

var _ KfgptaccountsentitiesModel = (*customKfgptaccountsentitiesModel)(nil)

type (
	// KfgptaccountsentitiesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customKfgptaccountsentitiesModel.
	KfgptaccountsentitiesModel interface {
		kfgptaccountsentitiesModel
	}

	customKfgptaccountsentitiesModel struct {
		*defaultKfgptaccountsentitiesModel
	}
)

// NewKfgptaccountsentitiesModel returns a model for the mongo.
func NewKfgptaccountsentitiesModel(url, db, collection string) KfgptaccountsentitiesModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customKfgptaccountsentitiesModel{
		defaultKfgptaccountsentitiesModel: newDefaultKfgptaccountsentitiesModel(conn),
	}
}
