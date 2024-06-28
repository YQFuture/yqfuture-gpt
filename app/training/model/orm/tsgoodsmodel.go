package orm

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TsGoodsModel = (*customTsGoodsModel)(nil)

type (
	// TsGoodsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTsGoodsModel.
	TsGoodsModel interface {
		tsGoodsModel
		withSession(session sqlx.Session) TsGoodsModel
	}

	customTsGoodsModel struct {
		*defaultTsGoodsModel
	}
)

// NewTsGoodsModel returns a model for the database table.
func NewTsGoodsModel(conn sqlx.SqlConn) TsGoodsModel {
	return &customTsGoodsModel{
		defaultTsGoodsModel: newTsGoodsModel(conn),
	}
}

func (m *customTsGoodsModel) withSession(session sqlx.Session) TsGoodsModel {
	return NewTsGoodsModel(sqlx.NewSqlConnFromSession(session))
}
