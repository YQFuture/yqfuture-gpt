package model

import "github.com/zeromicro/go-zero/core/stores/mon"

var _ ShoppresettingshoptitlesModel = (*customShoppresettingshoptitlesModel)(nil)

type (
	// ShoppresettingshoptitlesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShoppresettingshoptitlesModel.
	ShoppresettingshoptitlesModel interface {
		shoppresettingshoptitlesModel
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
