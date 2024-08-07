package orm

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BsShopCareHistoryModel = (*customBsShopCareHistoryModel)(nil)

type (
	// BsShopCareHistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBsShopCareHistoryModel.
	BsShopCareHistoryModel interface {
		bsShopCareHistoryModel
		withSession(session sqlx.Session) BsShopCareHistoryModel
		FindOrgShopCareData(ctx context.Context, shopId, startTime, endTime int64) (*OrgShopCareData, error)
	}

	customBsShopCareHistoryModel struct {
		*defaultBsShopCareHistoryModel
	}

	OrgShopCareData struct {
		CareTime  int64 `db:"care_time"`
		CareTimes int64 `db:"care_times"`
		UsedPower int64 `db:"used_power"`
	}
)

// NewBsShopCareHistoryModel returns a model for the database table.
func NewBsShopCareHistoryModel(conn sqlx.SqlConn) BsShopCareHistoryModel {
	return &customBsShopCareHistoryModel{
		defaultBsShopCareHistoryModel: newBsShopCareHistoryModel(conn),
	}
}

func (m *customBsShopCareHistoryModel) withSession(session sqlx.Session) BsShopCareHistoryModel {
	return NewBsShopCareHistoryModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBsShopCareHistoryModel) FindOrgShopCareData(ctx context.Context, shopId, startTime, endTime int64) (*OrgShopCareData, error) {
	var query string
	var resp OrgShopCareData
	var err error
	if startTime != 0 && endTime != 0 {
		query = fmt.Sprintf("select COALESCE(SUM(care_time), 0) AS care_time, COALESCE(SUM(used_power), 0) AS used_power, count(1) care_times from %s where `shop_id` = ? and `create_time` >= ? and `create_time` <= ?", m.table)
		err = m.conn.QueryRowCtx(ctx, &resp, query, shopId, startTime, endTime)
	} else {
		query = fmt.Sprintf("select COALESCE(SUM(care_time), 0) AS care_time, COALESCE(SUM(used_power), 0) AS used_power, count(1) care_times from %s where `shop_id` = ?", m.table)
		err = m.conn.QueryRowCtx(ctx, &resp, query, shopId)
	}
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
