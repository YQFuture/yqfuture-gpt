package orm

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BsDictTypeModel = (*customBsDictTypeModel)(nil)

type (
	// BsDictTypeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBsDictTypeModel.
	BsDictTypeModel interface {
		bsDictTypeModel
		withSession(session sqlx.Session) BsDictTypeModel
		FindOneByKey(ctx context.Context, key string) (*BsDictType, error)
	}

	customBsDictTypeModel struct {
		*defaultBsDictTypeModel
	}
)

// NewBsDictTypeModel returns a model for the database table.
func NewBsDictTypeModel(conn sqlx.SqlConn) BsDictTypeModel {
	return &customBsDictTypeModel{
		defaultBsDictTypeModel: newBsDictTypeModel(conn),
	}
}

func (m *customBsDictTypeModel) withSession(session sqlx.Session) BsDictTypeModel {
	return NewBsDictTypeModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBsDictTypeModel) FindOneByKey(ctx context.Context, key string) (*BsDictType, error) {
	query := fmt.Sprintf("select %s from %s where `key` = ? limit 1", bsDictTypeRows, m.table)
	var resp BsDictType
	err := m.conn.QueryRowCtx(ctx, &resp, query, key)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
