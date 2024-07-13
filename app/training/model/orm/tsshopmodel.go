package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
)

var _ TsShopModel = (*customTsShopModel)(nil)

type (
	// TsShopModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTsShopModel.
	TsShopModel interface {
		tsShopModel
		withSession(session sqlx.Session) TsShopModel
		FindList(ctx context.Context) (any, error)
		GetShopPageTotal(ctx context.Context, in *training.ShopPageListReq) (int, error)
		GetShopPageList(ctx context.Context, in *training.ShopPageListReq) ([]TsShop, error)
		JudgeShopFirst(ctx context.Context, userId int64, uuid string) (int, error)
		FindOneByUuidAndUserId(ctx context.Context, userId int64, uuid string) (*TsShop, error)
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

func (m *customTsShopModel) FindOneByUuidAndUserId(ctx context.Context, userId int64, uuid string) (*TsShop, error) {
	query := fmt.Sprintf("select %s from %s where user_id = ? AND uuid = ?", tsShopRows, m.table)
	var args []interface{}
	args = append(args, userId, uuid)
	var resp TsShop
	err := m.conn.QueryRowCtx(ctx, &resp, query, args...)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customTsShopModel) JudgeShopFirst(ctx context.Context, userId int64, uuid string) (int, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE user_id = ? AND uuid = ? AND training_times = 0", tsShopRows, m.table)
	var args []interface{}
	args = append(args, userId, uuid)
	var resp []*TsShop
	err := m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	//判断是否初次登录 0: 是 1: 否
	switch {
	case err == nil:
		if len(resp) > 0 {
			return 1, nil
		} else {
			return 0, nil
		}
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}

func (m *customTsShopModel) GetShopPageTotal(ctx context.Context, in *training.ShopPageListReq) (int, error) {
	// 动态构建WHERE子句
	var whereClauses []string
	var args []interface{}
	if in.UserId > 0 {
		whereClauses = append(whereClauses, "user_id = ?")
		args = append(args, in.UserId)
	}
	if in.Query != "" {
		whereClauses = append(whereClauses, "shop_name LIKE ?")
		args = append(args, "%"+in.Query+"%")
	}
	if in.PlatFormType > 0 {
		whereClauses = append(whereClauses, "platform_type = ?")
		args = append(args, in.PlatFormType)
	}
	if in.TrainingStatus > 0 {
		whereClauses = append(whereClauses, "training_status = ?")
		args = append(args, in.TrainingStatus)
	}
	if in.UpdateTime > 0 {
		whereClauses = append(whereClauses, "DATE(update_time) = DATE(FROM_UNIXTIME(?))")
		args = append(args, in.UpdateTime)
	}
	// 构建完整的SQL语句
	var whereStr string
	if len(whereClauses) > 0 {
		whereStr = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf("SELECT COUNT(1) FROM %s %s ORDER BY create_time DESC", m.table, whereStr)

	var count int
	err := m.conn.QueryRowCtx(ctx, &count, query, args...)

	switch {
	case err == nil:
		return count, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}

func (m *customTsShopModel) GetShopPageList(ctx context.Context, in *training.ShopPageListReq) ([]TsShop, error) {
	// 初始化偏移量和限制
	offset := (in.PageNum - 1) * in.PageSize
	limit := in.PageSize
	// 动态构建WHERE子句
	var whereClauses []string
	var args []interface{}
	if in.UserId > 0 {
		whereClauses = append(whereClauses, "user_id = ?")
		args = append(args, in.UserId)
	}
	if in.Query != "" {
		whereClauses = append(whereClauses, "shop_name LIKE ?")
		args = append(args, "%"+in.Query+"%")
	}
	if in.PlatFormType > 0 {
		whereClauses = append(whereClauses, "platform_type = ?")
		args = append(args, in.PlatFormType)
	}
	if in.TrainingStatus > 0 {
		whereClauses = append(whereClauses, "training_status = ?")
		args = append(args, in.TrainingStatus)
	}
	if in.UpdateTime > 0 {
		whereClauses = append(whereClauses, "DATE(update_time) = DATE(FROM_UNIXTIME(?))")
		args = append(args, in.UpdateTime)
	}
	// 构建完整的SQL语句
	var whereStr string
	if len(whereClauses) > 0 {
		whereStr = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY create_time DESC LIMIT ? OFFSET ?", tsShopRows, m.table, whereStr)
	args = append(args, limit, offset)

	var resp []TsShop
	err := m.conn.QueryRowsCtx(ctx, &resp, query, args...)

	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}
