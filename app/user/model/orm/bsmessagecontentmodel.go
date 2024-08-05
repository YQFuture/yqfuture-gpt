package orm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BsMessageContentModel = (*customBsMessageContentModel)(nil)

type (
	// BsMessageContentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBsMessageContentModel.
	BsMessageContentModel interface {
		bsMessageContentModel
		withSession(session sqlx.Session) BsMessageContentModel
		SessionInsert(ctx context.Context, data *BsMessageContent, session sqlx.Session) (sql.Result, error)
	}

	customBsMessageContentModel struct {
		*defaultBsMessageContentModel
	}
)

// NewBsMessageContentModel returns a model for the database table.
func NewBsMessageContentModel(conn sqlx.SqlConn) BsMessageContentModel {
	return &customBsMessageContentModel{
		defaultBsMessageContentModel: newBsMessageContentModel(conn),
	}
}

func (m *customBsMessageContentModel) withSession(session sqlx.Session) BsMessageContentModel {
	return NewBsMessageContentModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBsMessageContentModel) SessionInsert(ctx context.Context, data *BsMessageContent, session sqlx.Session) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, bsMessageContentRowsExpectAutoSet)
	ret, err := session.ExecCtx(ctx, query, data.Id, data.MessageType, data.MessageContentType, data.MessageContent, data.CreateBy, data.UpdateBy)
	return ret, err
}
