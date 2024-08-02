package orm

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BsShopModel = (*customBsShopModel)(nil)

type (
	// BsShopModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBsShopModel.
	BsShopModel interface {
		bsShopModel
		withSession(session sqlx.Session) BsShopModel
		FindListByOrgId(ctx context.Context, orgId int64) (*[]*BsShop, error)
		FindOrgTotalMonthUsedPower(ctx context.Context, orgId int64) (int64, error)
		UpdateShopPower(ctx context.Context, shopId, power int64) error
		UpdateShopPowerAvg(ctx context.Context, orgId, power int64) error
	}

	customBsShopModel struct {
		*defaultBsShopModel
	}
)

// NewBsShopModel returns a model for the database table.
func NewBsShopModel(conn sqlx.SqlConn) BsShopModel {
	return &customBsShopModel{
		defaultBsShopModel: newBsShopModel(conn),
	}
}

func (m *customBsShopModel) withSession(session sqlx.Session) BsShopModel {
	return NewBsShopModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBsShopModel) FindListByOrgId(ctx context.Context, orgId int64) (*[]*BsShop, error) {
	query := fmt.Sprintf("select %s from %s where `org_id` = ?", bsShopRows, m.table)
	var resp []*BsShop
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

func (m *customBsShopModel) FindOrgTotalMonthUsedPower(ctx context.Context, orgId int64) (int64, error) {
	query := fmt.Sprintf("select count(month_used_power) from %s where `org_id` = ?", m.table)
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

func (m *customBsShopModel) UpdateShopPower(ctx context.Context, shopId, power int64) error {
	query := fmt.Sprintf("update %s set month_power_limit = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, power, shopId)
	return err
}

func (m *customBsShopModel) UpdateShopPowerAvg(ctx context.Context, orgId, power int64) error {
	query := fmt.Sprintf("update %s set month_power_limit = ? where org_id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, power, orgId)
	return err
}
