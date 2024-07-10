package orm

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TsShopLogModel = (*customTsShopLogModel)(nil)

type (
	// TsShopLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTsShopLogModel.
	TsShopLogModel interface {
		tsShopLogModel
		withSession(session sqlx.Session) TsShopLogModel
	}

	customTsShopLogModel struct {
		*defaultTsShopLogModel
	}
)

// NewTsShopLogModel returns a model for the database table.
func NewTsShopLogModel(conn sqlx.SqlConn) TsShopLogModel {
	return &customTsShopLogModel{
		defaultTsShopLogModel: newTsShopLogModel(conn),
	}
}

func (m *customTsShopLogModel) withSession(session sqlx.Session) TsShopLogModel {
	return NewTsShopLogModel(sqlx.NewSqlConnFromSession(session))
}
