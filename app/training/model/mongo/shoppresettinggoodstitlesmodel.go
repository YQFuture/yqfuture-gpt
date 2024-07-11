package model

import "github.com/zeromicro/go-zero/core/stores/mon"

var _ ShoppresettinggoodstitlesModel = (*customShoppresettinggoodstitlesModel)(nil)

type (
	// ShoppresettinggoodstitlesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShoppresettinggoodstitlesModel.
	ShoppresettinggoodstitlesModel interface {
		shoppresettinggoodstitlesModel
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
