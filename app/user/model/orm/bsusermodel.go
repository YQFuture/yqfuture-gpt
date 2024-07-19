package orm

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BsUserModel = (*customBsUserModel)(nil)

type (
	// BsUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBsUserModel.
	BsUserModel interface {
		bsUserModel
		withSession(session sqlx.Session) BsUserModel
		FindOneByPhone(ctx context.Context, phone string) (*BsUser, error)
		FindOneByOpenId(ctx context.Context, openid string) (*BsUser, error)
		BindPhone(ctx context.Context, phone string, userId int64) error
		BindOpenId(ctx context.Context, openid string, userId int64) error
	}

	customBsUserModel struct {
		*defaultBsUserModel
	}
)

// NewBsUserModel returns a model for the database table.
func NewBsUserModel(conn sqlx.SqlConn) BsUserModel {
	return &customBsUserModel{
		defaultBsUserModel: newBsUserModel(conn),
	}
}

func (m *customBsUserModel) withSession(session sqlx.Session) BsUserModel {
	return NewBsUserModel(sqlx.NewSqlConnFromSession(session))
}

// FindOneByPhone 根据手机号查询用户
func (m *customBsUserModel) FindOneByPhone(ctx context.Context, phone string) (*BsUser, error) {
	query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", bsUserRows, m.table)
	var resp BsUser
	err := m.conn.QueryRowCtx(ctx, &resp, query, phone)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

// FindOneByOpenId 根据OpenId查询用户
func (m *customBsUserModel) FindOneByOpenId(ctx context.Context, openid string) (*BsUser, error) {
	query := fmt.Sprintf("select %s from %s where `openid` = ? limit 1", bsUserRows, m.table)
	var resp BsUser
	err := m.conn.QueryRowCtx(ctx, &resp, query, openid)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customBsUserModel) BindPhone(ctx context.Context, phone string, userId int64) error {
	query := fmt.Sprintf("update %s set `phone` = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, phone, userId)
	return err
}

func (m *customBsUserModel) BindOpenId(ctx context.Context, openid string, userId int64) error {
	query := fmt.Sprintf("update %s set `openid` = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, openid, userId)
	return err
}
