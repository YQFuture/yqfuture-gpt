package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
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
		TransactCtx(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		SessionInsert(ctx context.Context, data *BsUser, session sqlx.Session) (sql.Result, error)
		ChangeOrg(ctx context.Context, orgId, userId int64) error
		UpdateHeadImg(ctx context.Context, headImg string, userId int64) error
		UpdateNickName(ctx context.Context, nickName string, userId int64) error
		FindListByPhone(ctx context.Context, queryString string) (*[]*BsUser, error)
		FindPageListByOrgId(ctx context.Context, orgId, pageNum, pageSize int64, queryString string) (*[]*OrgUser, error)
		FindPageTotalByOrgId(ctx context.Context, orgId, pageNum, pageSize int64, queryString string) (int64, error)
		FindListByOrgId(ctx context.Context, orgId int64) (*[]*BsUser, error)
	}

	customBsUserModel struct {
		*defaultBsUserModel
	}

	OrgUser struct {
		Id              int64          `db:"id"`                // 用户ID
		NowOrgId        int64          `db:"now_org_id"`        // 当前组织ID
		UserName        sql.NullString `db:"user_name"`         // 用户名
		NickName        sql.NullString `db:"nick_name"`         // 用户昵称
		HeadImg         sql.NullString `db:"head_img"`          // 头像地址
		Phone           sql.NullString `db:"phone"`             // 手机号码
		Password        sql.NullString `db:"password"`          // 密码
		Openid          sql.NullString `db:"openid"`            // openid是微信用户在不同类型的产品的身份ID
		Unionid         sql.NullString `db:"unionid"`           // unionid是微信用户在同一个开放平台下的产品的身份ID
		CreateTime      time.Time      `db:"create_time"`       // 创建时间
		UpdateTime      time.Time      `db:"update_time"`       // 修改时间
		CreateBy        int64          `db:"create_by"`         // 创建人
		UpdateBy        int64          `db:"update_by"`         // 修改人
		Status          int64          `db:"status"`            // 状态 0: 暂停 1: 启用
		MonthPowerLimit int64          `db:"month_power_limit"` // 当月算力上限
		MonthUsedPower  int64          `db:"month_used_power"`  // 当月已用算力
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

func (m *customBsUserModel) TransactCtx(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customBsUserModel) SessionInsert(ctx context.Context, data *BsUser, session sqlx.Session) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, bsUserRowsExpectAutoSet)
	ret, err := session.ExecCtx(ctx, query, data.Id, data.NowOrgId, data.UserName, data.NickName, data.HeadImg, data.Phone, data.Password, data.Openid, data.Unionid, data.CreateBy, data.UpdateBy)
	return ret, err
}

func (m *customBsUserModel) ChangeOrg(ctx context.Context, orgId, userId int64) error {
	query := fmt.Sprintf("update %s set `now_org_id` = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, orgId, userId)
	return err
}

func (m *customBsUserModel) UpdateHeadImg(ctx context.Context, headImg string, userId int64) error {
	query := fmt.Sprintf("update %s set `head_img` = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, headImg, userId)
	return err
}

func (m *customBsUserModel) UpdateNickName(ctx context.Context, nickName string, userId int64) error {
	query := fmt.Sprintf("update %s set `nick_name` = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, nickName, userId)
	return err
}

func (m *customBsUserModel) FindListByPhone(ctx context.Context, phone string) (*[]*BsUser, error) {
	query := fmt.Sprintf("SELECT * FROM bs_user WHERE `phone` = ?")
	var resp []*BsUser
	err := m.conn.QueryRowsCtx(ctx, &resp, query, phone)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customBsUserModel) FindPageListByOrgId(ctx context.Context, orgId, pageNum, pageSize int64, queryString string) (*[]*OrgUser, error) {
	// 初始化偏移量和限制
	offset := (pageNum - 1) * pageSize
	limit := pageSize
	var query string
	var resp []*OrgUser
	var err error
	if queryString != "" {
		query = fmt.Sprintf("SELECT u.*,uo.status,uo.month_power_limit,uo.month_used_power FROM bs_user_org uo LEFT JOIN bs_user u ON uo.user_id = u.id WHERE uo.org_id = ? AND (u.nick_name LIKE ? OR u.phone LIKE ?) ORDER BY uo.create_time ASC LIMIT ? OFFSET ?")
		err = m.conn.QueryRowsCtx(ctx, &resp, query, orgId, "%"+queryString+"%", "%"+queryString+"%", limit, offset)
	} else {
		query = fmt.Sprintf("SELECT u.*,uo.status,uo.month_power_limit,uo.month_used_power FROM bs_user_org uo LEFT JOIN bs_user u ON uo.user_id = u.id WHERE uo.org_id = ? ORDER BY uo.create_time ASC LIMIT ? OFFSET ?")
		err = m.conn.QueryRowsCtx(ctx, &resp, query, orgId, limit, offset)
	}
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customBsUserModel) FindPageTotalByOrgId(ctx context.Context, orgId, pageNum, pageSize int64, queryString string) (int64, error) {
	// 初始化偏移量和限制
	offset := (pageNum - 1) * pageSize
	limit := pageSize
	var query string
	var resp int64
	var err error
	if queryString != "" {
		query = fmt.Sprintf("SELECT COUNT(1) FROM bs_user_org uo LEFT JOIN bs_user u ON uo.user_id = u.id WHERE uo.org_id = ? AND (u.nick_name LIKE ? OR u.phone LIKE ?) ORDER BY uo.create_time ASC LIMIT ? OFFSET ?")
		err = m.conn.QueryRowCtx(ctx, &resp, query, orgId, "%"+queryString+"%", "%"+queryString+"%", limit, offset)
	} else {
		query = fmt.Sprintf("SELECT COUNT(1) FROM bs_user_org uo LEFT JOIN bs_user u ON uo.user_id = u.id WHERE uo.org_id = ? ORDER BY uo.create_time ASC LIMIT ? OFFSET ?")
		err = m.conn.QueryRowCtx(ctx, &resp, query, orgId, limit, offset)
	}
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}

}

func (m *customBsUserModel) FindListByOrgId(ctx context.Context, orgId int64) (*[]*BsUser, error) {
	query := fmt.Sprintf("SELECT u.* FROM bs_user_org uo LEFT JOIN bs_user u ON uo.user_id = u.id WHERE uo.org_id = ?")
	var resp []*BsUser
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
