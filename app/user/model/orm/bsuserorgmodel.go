package orm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BsUserOrgModel = (*customBsUserOrgModel)(nil)

type (
	// BsUserOrgModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBsUserOrgModel.
	BsUserOrgModel interface {
		bsUserOrgModel
		withSession(session sqlx.Session) BsUserOrgModel
		SessionInsert(ctx context.Context, data *BsUserOrg, session sqlx.Session) (sql.Result, error)
	}

	customBsUserOrgModel struct {
		*defaultBsUserOrgModel
	}
)

// NewBsUserOrgModel returns a model for the database table.
func NewBsUserOrgModel(conn sqlx.SqlConn) BsUserOrgModel {
	return &customBsUserOrgModel{
		defaultBsUserOrgModel: newBsUserOrgModel(conn),
	}
}

func (m *customBsUserOrgModel) withSession(session sqlx.Session) BsUserOrgModel {
	return NewBsUserOrgModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultBsUserOrgModel) SessionInsert(ctx context.Context, data *BsUserOrg, session sqlx.Session) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, bsUserOrgRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.OrgId)
	return ret, err
}
