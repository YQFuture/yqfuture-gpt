package orm

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TsShopModel = (*customTsShopModel)(nil)

type (
	// TsShopModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTsShopModel.
	TsShopModel interface {
		tsShopModel
		withSession(session sqlx.Session) TsShopModel
		FindList(ctx context.Context) (any, error)
	}

	customTsShopModel struct {
		*defaultTsShopModel
	}
)

// NewTsShopModel returns a model for the database table.
func NewTsShopModel(conn sqlx.SqlConn) TsShopModel {
	return &customTsShopModel{
		defaultTsShopModel: newTsShopModel(conn),
	}
}

func (m *customTsShopModel) withSession(session sqlx.Session) TsShopModel {
	return NewTsShopModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customTsShopModel) FindList(ctx context.Context) (any, error) {
	query := fmt.Sprintf("select %s from %s", tsShopRows, m.table)
	var resp []*TsShop
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
