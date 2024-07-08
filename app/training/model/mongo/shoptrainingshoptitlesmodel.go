package model

import "github.com/zeromicro/go-zero/core/stores/mon"

var _ ShoptrainingshoptitlesModel = (*customShoptrainingshoptitlesModel)(nil)

type (
	// ShoptrainingshoptitlesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShoptrainingshoptitlesModel.
	ShoptrainingshoptitlesModel interface {
		shoptrainingshoptitlesModel
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
