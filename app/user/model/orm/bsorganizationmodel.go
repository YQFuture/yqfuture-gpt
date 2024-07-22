package orm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BsOrganizationModel = (*customBsOrganizationModel)(nil)

type (
	// BsOrganizationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBsOrganizationModel.
	BsOrganizationModel interface {
		bsOrganizationModel
		withSession(session sqlx.Session) BsOrganizationModel
		SessionInsert(ctx context.Context, data *BsOrganization, session sqlx.Session) (sql.Result, error)
	}

	customBsOrganizationModel struct {
		*defaultBsOrganizationModel
	}
)

// NewBsOrganizationModel returns a model for the database table.
func NewBsOrganizationModel(conn sqlx.SqlConn) BsOrganizationModel {
	return &customBsOrganizationModel{
		defaultBsOrganizationModel: newBsOrganizationModel(conn),
	}
}

func (m *customBsOrganizationModel) withSession(session sqlx.Session) BsOrganizationModel {
	return NewBsOrganizationModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultBsOrganizationModel) SessionInsert(ctx context.Context, data *BsOrganization, session sqlx.Session) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, bsOrganizationRowsExpectAutoSet)
	ret, err := session.ExecCtx(ctx, query, data.Id, data.OwnerId, data.OrgName, data.BundleType, data.CreateBy, data.UpdateBy)
	return ret, err
}
