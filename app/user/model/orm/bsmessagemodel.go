package orm

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BsMessageModel = (*customBsMessageModel)(nil)

type (
	// BsMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBsMessageModel.
	BsMessageModel interface {
		bsMessageModel
		withSession(session sqlx.Session) BsMessageModel
		SyncNotice(ctx context.Context, userId int64) error
		FindUnreadCount(ctx context.Context, userId, nowOrgId int64) (int64, error)
	}

	customBsMessageModel struct {
		*defaultBsMessageModel
	}
)

// NewBsMessageModel returns a model for the database table.
func NewBsMessageModel(conn sqlx.SqlConn) BsMessageModel {
	return &customBsMessageModel{
		defaultBsMessageModel: newBsMessageModel(conn),
	}
}

func (m *customBsMessageModel) withSession(session sqlx.Session) BsMessageModel {
	return NewBsMessageModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultBsMessageModel) SyncNotice(ctx context.Context, userId int64) error {
	query := fmt.Sprintf("INSERT INTO bs_message ( user_id, org_id, message_type, content_id, read_flag, ignore_flag, create_time, update_time, create_by, update_by ) SELECT ? AS user_id, NULL AS org_id, message_type, id AS content_id, 0 AS read_flag, 0 AS ignore_flag, NOW() AS create_time, NOW() AS update_time, create_by, update_by FROM   bs_message_content WHERE id > ( SELECT COALESCE ( MAX( content_id ), 0 ) FROM bs_message )")
	_, err := m.conn.ExecCtx(ctx, query, userId)
	return err
}

func (m *defaultBsMessageModel) FindUnreadCount(ctx context.Context, userId, nowOrgId int64) (int64, error) {
	query := fmt.Sprintf("select count(1) from %s where `user_id` = ? and `read_flag` = 0 and (org_id = ? or org_id is null)", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, nowOrgId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, ErrNotFound
	default:
		return 0, err
	}
}
