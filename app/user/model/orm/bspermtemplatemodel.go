package orm

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BsPermTemplateModel = (*customBsPermTemplateModel)(nil)

type (
	// BsPermTemplateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBsPermTemplateModel.
	BsPermTemplateModel interface {
		bsPermTemplateModel
		withSession(session sqlx.Session) BsPermTemplateModel
		FindListByBundleType(ctx context.Context, bundleType int) (*[]*BsPermTemplate, error)
	}

	customBsPermTemplateModel struct {
		*defaultBsPermTemplateModel
	}
)

// NewBsPermTemplateModel returns a model for the database table.
func NewBsPermTemplateModel(conn sqlx.SqlConn) BsPermTemplateModel {
	return &customBsPermTemplateModel{
		defaultBsPermTemplateModel: newBsPermTemplateModel(conn),
	}
}

func (m *customBsPermTemplateModel) withSession(session sqlx.Session) BsPermTemplateModel {
	return NewBsPermTemplateModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBsPermTemplateModel) FindListByBundleType(ctx context.Context, bundleType int) (*[]*BsPermTemplate, error) {
	query := fmt.Sprintf("select %s from %s where `bundle_type` = ? and `enable_flag` = 1", bsPermTemplateRows, m.table)
	var resp []*BsPermTemplate
	err := m.conn.QueryRowsCtx(ctx, &resp, query, bundleType)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
