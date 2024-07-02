package orm

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
)

var _ TsGoodsModel = (*customTsGoodsModel)(nil)

type (
	// TsGoodsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTsGoodsModel.
	TsGoodsModel interface {
		tsGoodsModel
		withSession(session sqlx.Session) TsGoodsModel
		GetGoodsPageList(ctx context.Context, in *training.GoodsPageListReq) (*[]*TsGoods, error)
		GetGoodsPageTotal(ctx context.Context, in *training.GoodsPageListReq) (int, error)
		EnableGoods(ctx context.Context, in *training.GoodsTrainingReq) error
		UnEnableGoods(ctx context.Context, in *training.GoodsTrainingReq) error
		FindEnabledListByShopId(ctx context.Context, in int64) (*[]*TsGoods, error)
	}

	customTsGoodsModel struct {
		*defaultTsGoodsModel
	}
)

// NewTsGoodsModel returns a model for the database table.
func NewTsGoodsModel(conn sqlx.SqlConn) TsGoodsModel {
	return &customTsGoodsModel{
		defaultTsGoodsModel: newTsGoodsModel(conn),
	}
}

func (m *customTsGoodsModel) withSession(session sqlx.Session) TsGoodsModel {
	return NewTsGoodsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customTsGoodsModel) FindEnabledListByShopId(ctx context.Context, in int64) (*[]*TsGoods, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE enabled = 2 AND shop_id = ?", tsGoodsRows, m.table)
	var resp []*TsGoods
	err := m.conn.QueryRowsCtx(ctx, &resp, query, in)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customTsGoodsModel) GetGoodsPageTotal(ctx context.Context, in *training.GoodsPageListReq) (int, error) {
	// 动态构建WHERE子句
	var whereClauses []string
	var args []interface{}

	if in.ShopId > 0 {
		whereClauses = append(whereClauses, "shop_id = ?")
		args = append(args, in.ShopId)
	}
	if in.Query != "" {
		whereClauses = append(whereClauses, "goods_name LIKE ?")
		args = append(args, "%"+in.Query+"%")
	}
	if in.Enabled > 0 {
		whereClauses = append(whereClauses, "enabled = ?")
		args = append(args, in.Enabled)
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

	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY create_time DESC", tsGoodsRows, m.table, whereStr)

	var resp []*TsGoods
	err := m.conn.QueryRowsCtx(ctx, &resp, query, args...)

	switch err {
	case nil:
		return len(resp), nil
	case sqlx.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

func (m *customTsGoodsModel) GetGoodsPageList(ctx context.Context, in *training.GoodsPageListReq) (*[]*TsGoods, error) {
	// 初始化偏移量和限制
	offset := (in.PageNum - 1) * in.PageSize
	limit := in.PageSize

	// 动态构建WHERE子句
	var whereClauses []string
	var args []interface{}

	if in.ShopId > 0 {
		whereClauses = append(whereClauses, "shop_id = ?")
		args = append(args, in.ShopId)
	}
	if in.Query != "" {
		whereClauses = append(whereClauses, "goods_name LIKE ?")
		args = append(args, "%"+in.Query+"%")
	}
	if in.Enabled > 0 {
		whereClauses = append(whereClauses, "enabled = ?")
		args = append(args, in.Enabled)
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

	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY create_time DESC LIMIT ? OFFSET ?", tsGoodsRows, m.table, whereStr)
	args = append(args, limit, offset)

	var resp []*TsGoods
	err := m.conn.QueryRowsCtx(ctx, &resp, query, args...)

	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customTsGoodsModel) EnableGoods(ctx context.Context, in *training.GoodsTrainingReq) error {
	_, err := m.conn.ExecCtx(ctx, "UPDATE ts_goods SET enabled = 2 WHERE id = ?", in.GoodsId)
	if err != nil {
		return err
	}
	return nil
}

func (m *customTsGoodsModel) UnEnableGoods(ctx context.Context, in *training.GoodsTrainingReq) error {
	_, err := m.conn.ExecCtx(ctx, "UPDATE ts_goods SET enabled = 1 WHERE id = ?", in.GoodsId)
	if err != nil {
		return err
	}
	return nil
}
