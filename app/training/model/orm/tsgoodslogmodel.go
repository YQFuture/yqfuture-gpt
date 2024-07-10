package orm

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TsGoodsLogModel = (*customTsGoodsLogModel)(nil)

type (
	// TsGoodsLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTsGoodsLogModel.
	TsGoodsLogModel interface {
		tsGoodsLogModel
		withSession(session sqlx.Session) TsGoodsLogModel
	}

	customTsGoodsLogModel struct {
		*defaultTsGoodsLogModel
	}
)

// NewTsGoodsLogModel returns a model for the database table.
func NewTsGoodsLogModel(conn sqlx.SqlConn) TsGoodsLogModel {
	return &customTsGoodsLogModel{
		defaultTsGoodsLogModel: newTsGoodsLogModel(conn),
	}
}

func (m *customTsGoodsLogModel) withSession(session sqlx.Session) TsGoodsLogModel {
	return NewTsGoodsLogModel(sqlx.NewSqlConnFromSession(session))
}
