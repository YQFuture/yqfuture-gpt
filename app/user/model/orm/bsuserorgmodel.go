package orm

import (
	"context"
	"database/sql"
	"errors"
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
		FindUserOrgCount(ctx context.Context, userId int64) (int64, error)
		FindUserOrgByUserIdAndOrgId(ctx context.Context, userId, orgId int64) (*BsUserOrg, error)
		FindOrgUserCount(ctx context.Context, orgId int64) (int64, error)
		DeleteByUserIdAndOrgId(ctx context.Context, userId, orgId int64) error
		ChangeStatusByUserIdAndOrgId(ctx context.Context, userId, orgId, status int64) error
		FindOrgTotalGivePower(ctx context.Context, orgId int64) (int64, error)
		UpdateUserPower(ctx context.Context, userId, orgId, power int64) error
		FindListByOrgId(ctx context.Context, orgId int64) (*[]*BsUserOrg, error)
		UpdateUserPowerAvg(ctx context.Context, orgId, power int64) error
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

func (m *customBsUserOrgModel) SessionInsert(ctx context.Context, data *BsUserOrg, session sqlx.Session) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, bsUserOrgRowsExpectAutoSet)
	ret, err := session.ExecCtx(ctx, query, data.UserId, data.OrgId, data.Status, data.MonthPowerLimit, data.MonthUsedPower, data.CreateBy, data.UpdateBy)
	return ret, err
}

func (m *customBsUserOrgModel) FindUserOrgCount(ctx context.Context, userId int64) (int64, error) {
	query := fmt.Sprintf("select count(1) from %s where `user_id` = ?", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}

func (m *customBsUserOrgModel) FindUserOrgByUserIdAndOrgId(ctx context.Context, userId, orgId int64) (*BsUserOrg, error) {
	query := fmt.Sprintf("select * from %s where `user_id` = ? and `org_id` = ?", m.table)
	var resp BsUserOrg
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

func (m *customBsUserOrgModel) FindOrgUserCount(ctx context.Context, orgId int64) (int64, error) {
	query := fmt.Sprintf("select count(1) from %s where `org_id` = ?", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query, orgId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}

func (m *customBsUserOrgModel) DeleteByUserIdAndOrgId(ctx context.Context, userId, orgId int64) error {
	query := fmt.Sprintf("delete from %s where `user_id` = ? and org_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId, orgId)
	return err
}

func (m *customBsUserOrgModel) ChangeStatusByUserIdAndOrgId(ctx context.Context, userId, orgId, status int64) error {
	query := fmt.Sprintf("update %s set status = ? where `user_id` = ? and org_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, status, userId, orgId)
	return err
}

func (m *customBsUserOrgModel) FindOrgTotalGivePower(ctx context.Context, orgId int64) (int64, error) {
	query := fmt.Sprintf("select SUM(month_power_limit) from %s where `org_id` = ?", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query, orgId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}

func (m *customBsUserOrgModel) UpdateUserPower(ctx context.Context, userId, orgId, power int64) error {
	query := fmt.Sprintf("update %s set month_power_limit = ? where `user_id` = ? and org_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, power, userId, orgId)
	return err
}

func (m *customBsUserOrgModel) UpdateUserPowerAvg(ctx context.Context, orgId, power int64) error {
	query := fmt.Sprintf("update %s set month_power_limit = ? where org_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, power, orgId)
	return err
}

func (m *customBsUserOrgModel) FindListByOrgId(ctx context.Context, orgId int64) (*[]*BsUserOrg, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE `oeg_id` = ?", m.table)
	var resp []*BsUserOrg
	err := m.conn.QueryRowsCtx(ctx, &resp, query, orgId)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}
