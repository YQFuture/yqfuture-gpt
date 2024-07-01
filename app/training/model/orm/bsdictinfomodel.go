package orm

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BsDictInfoModel = (*customBsDictInfoModel)(nil)

type (
	// BsDictInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBsDictInfoModel.
	BsDictInfoModel interface {
		bsDictInfoModel
		withSession(session sqlx.Session) BsDictInfoModel
		FindListByTypeId(ctx context.Context, typeId int64) (*[]*BsDictInfo, error)
	}

	customBsDictInfoModel struct {
		*defaultBsDictInfoModel
	}
)

// NewBsDictInfoModel returns a model for the database table.
func NewBsDictInfoModel(conn sqlx.SqlConn) BsDictInfoModel {
	return &customBsDictInfoModel{
		defaultBsDictInfoModel: newBsDictInfoModel(conn),
	}
}

func (m *customBsDictInfoModel) withSession(session sqlx.Session) BsDictInfoModel {
	return NewBsDictInfoModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBsDictInfoModel) FindListByTypeId(ctx context.Context, typeId int64) (*[]*BsDictInfo, error) {
	query := fmt.Sprintf("select %s from %s where `typeId` = ?", bsDictInfoRows, m.table)
	var resp []*BsDictInfo
	err := m.conn.QueryRowsCtx(ctx, &resp, query, typeId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
