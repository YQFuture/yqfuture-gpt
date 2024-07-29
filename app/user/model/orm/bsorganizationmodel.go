package orm

import (
	"context"
	"database/sql"
	"errors"
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
		FindListByUserId(ctx context.Context, userId int64) (*[]*BsOrganization, error)
		FindOneByIdAndUserId(ctx context.Context, id, userId int64) (*BsOrganization, error)
		FindOneByName(ctx context.Context, orgName string) (*BsOrganization, error)
		UpdateOrgName(ctx context.Context, orgName string, orgId int64) error
		FindListByNameOrOwnerPhone(ctx context.Context, queryString string) (*[]*BsOrganization, error)
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

func (m *customBsOrganizationModel) SessionInsert(ctx context.Context, data *BsOrganization, session sqlx.Session) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, bsOrganizationRowsExpectAutoSet)
	ret, err := session.ExecCtx(ctx, query, data.Id, data.OwnerId, data.OrgName, data.BundleType, data.CreateBy, data.UpdateBy)
	return ret, err
}

func (m *customBsOrganizationModel) FindListByUserId(ctx context.Context, userId int64) (*[]*BsOrganization, error) {
	query := fmt.Sprintf("SELECT o.* FROM bs_organization o INNER JOIN bs_user_org uo ON uo.org_id = o.id WHERE uo.user_id = ?")
	var resp []*BsOrganization
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customBsOrganizationModel) FindOneByIdAndUserId(ctx context.Context, orgId, userId int64) (*BsOrganization, error) {
	query := fmt.Sprintf("SELECT o.* FROM bs_organization o INNER JOIN ( SELECT org_id FROM bs_user_org WHERE user_id = ? ) uo ON o.id = uo.org_id WHERE o.id = ?")
	var resp BsOrganization
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, orgId)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customBsOrganizationModel) FindOneByName(ctx context.Context, orgName string) (*BsOrganization, error) {
	query := fmt.Sprintf("select %s from %s where `org_name` = ? limit 1", bsOrganizationRows, m.table)
	var resp BsOrganization
	err := m.conn.QueryRowCtx(ctx, &resp, query, orgName)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customBsOrganizationModel) UpdateOrgName(ctx context.Context, orgName string, orgId int64) error {
	query := fmt.Sprintf("update %s set org_name = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, orgName, orgId)
	return err
}

func (m *customBsOrganizationModel) FindListByNameOrOwnerPhone(ctx context.Context, queryString string) (*[]*BsOrganization, error) {
	query := fmt.Sprintf("SELECT o.* FROM bs_organization o LEFT JOIN bs_user u ON o.owner_id = u.id WHERE o.org_name = ? OR u.phone = ?")
	var resp []*BsOrganization
	err := m.conn.QueryRowsCtx(ctx, &resp, query, queryString, queryString)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}
