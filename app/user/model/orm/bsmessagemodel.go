package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
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
		FindMessageList(ctx context.Context, userId, nowOrgId, messageId, timeVector int64) (*[]*BsMessageInfo, error)
		SetMessageRead(ctx context.Context, userId, nowOrgId int64) error
		IgnoreMessage(ctx context.Context, messageId int64) error
		TransactCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		SessionInsert(ctx context.Context, data *BsMessage, session sqlx.Session) (sql.Result, error)
	}

	customBsMessageModel struct {
		*defaultBsMessageModel
	}

	BsMessageInfo struct {
		Id                 int64     `db:"id"`                   // 消息ID
		MessageContentType int64     `db:"message_content_type"` // 消息内容类型 0: 文字 1: 图文 2: 图片 3: 邀请加入组织 4: 申请加入组织 5: 平台掉线 6: 算力不足 7: 转接失败
		MessageContent     string    `db:"message_content"`      // 消息内容
		DealFlag           int64     `db:"deal_flag"`            // 处理标记 0: 未处理 1: 已处理
		IgnoreFlag         int64     `db:"ignore_flag"`          // 忽略标记 0: 未忽略 1: 已忽略
		CreateTime         time.Time `db:"create_time"`          // 创建时间
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

func (m *customBsMessageModel) SyncNotice(ctx context.Context, userId int64) error {
	query := fmt.Sprintf("INSERT INTO bs_message ( user_id, org_id, message_type, content_id, read_flag, deal_flag, ignore_flag, create_time, update_time, create_by, update_by ) SELECT ? AS user_id, 0 AS org_id, message_type, id AS content_id, 0 AS read_flag, 0 AS deal_flag, 0 AS ignore_flag, NOW() AS create_time, NOW() AS update_time, create_by, update_by FROM   bs_message_content WHERE   message_type = 1   AND id > (   SELECT COALESCE     ( MAX( content_id ), 0 )   FROM     bs_message   WHERE   message_type = 1 AND `user_id` = ?)")
	_, err := m.conn.ExecCtx(ctx, query, userId, userId)
	return err
}

func (m *customBsMessageModel) FindUnreadCount(ctx context.Context, userId, nowOrgId int64) (int64, error) {
	query := fmt.Sprintf("select count(1) from %s where `user_id` = ? and `read_flag` = 0 and (org_id = ? or org_id = 0)", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, nowOrgId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}

func (m *customBsMessageModel) FindMessageList(ctx context.Context, userId, nowOrgId, messageId, timeVector int64) (*[]*BsMessageInfo, error) {
	var resp []*BsMessageInfo
	var err error
	if messageId != 0 {
		if timeVector == 0 {
			query := fmt.Sprintf("SELECT m.id, c.message_content_type, c.message_content, m.deal_flag, m.ignore_flag, c.create_time FROM bs_message m LEFT JOIN bs_message_content c ON m.content_id = c.id WHERE m.id > ? AND m.user_id = ? AND (m.org_id = ? OR m.org_id = 0) ORDER BY m.id DESC LIMIT 5")
			err = m.conn.QueryRowsCtx(ctx, &resp, query, messageId, userId, nowOrgId)
		} else {
			query := fmt.Sprintf("SELECT m.id, c.message_content_type, c.message_content, m.deal_flag, m.ignore_flag, c.create_time FROM bs_message m LEFT JOIN bs_message_content c ON m.content_id = c.id WHERE m.id < ? AND m.user_id = ? AND (m.org_id = ? OR m.org_id = 0) ORDER BY m.id DESC LIMIT 5")
			err = m.conn.QueryRowsCtx(ctx, &resp, query, messageId, userId, nowOrgId)
		}
	} else {
		query := fmt.Sprintf("SELECT m.id, c.message_content_type, c.message_content, m.deal_flag, m.ignore_flag, c.create_time FROM bs_message m LEFT JOIN bs_message_content c ON m.content_id = c.id WHERE m.user_id = ? AND (m.org_id = ? OR m.org_id = 0) ORDER BY m.id DESC LIMIT 5")
		err = m.conn.QueryRowsCtx(ctx, &resp, query, userId, nowOrgId)
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

func (m *customBsMessageModel) SetMessageRead(ctx context.Context, userId, nowOrgId int64) error {
	query := fmt.Sprintf("update %s set read_flag = 1 where user_id = ? and (org_id = ? OR org_id = 0) and `read_flag` = 0", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId, nowOrgId)
	return err
}

func (m *customBsMessageModel) IgnoreMessage(ctx context.Context, messageId int64) error {
	query := fmt.Sprintf("update %s set ignore_flag = 1 where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, messageId)
	return err
}

func (m *customBsMessageModel) TransactCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customBsMessageModel) SessionInsert(ctx context.Context, data *BsMessage, session sqlx.Session) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, bsMessageRowsExpectAutoSet)
	ret, err := session.ExecCtx(ctx, query, data.UserId, data.OrgId, data.MessageType, data.ContentId, data.ReadFlag, data.DealFlag, data.IgnoreFlag, data.CreateBy, data.UpdateBy)
	return ret, err
}
