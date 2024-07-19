package userlogic

import (
	"context"
	"database/sql"
	"time"
	"yufuture-gpt/app/user/model/orm"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindPhoneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindPhoneLogic {
	return &BindPhoneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BindPhone 绑定手机号
func (l *BindPhoneLogic) BindPhone(in *user.BindPhoneReq) (*user.BindPhoneResp, error) {
	bsUser, err := l.svcCtx.BsUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		return nil, err
	}
	// 如果手机号未注册 创建新用户 并返回当前用户ID
	if bsUser == nil {
		now := time.Now()
		id := l.svcCtx.SnowFlakeNode.Generate().Int64()
		newBsUser := &orm.BsUser{
			Id: id,
			Phone: sql.NullString{
				String: in.Phone,
				Valid:  true,
			},
			UserName: sql.NullString{
				String: in.Phone,
				Valid:  true,
			},
			NickName: sql.NullString{
				String: in.Phone,
				Valid:  true,
			},
			HeadImg: sql.NullString{
				String: "https://yqfuture.com/_nuxt/favicon.28e9763f.png",
				Valid:  true,
			},
			Openid: sql.NullString{
				String: in.Openid,
				Valid:  true,
			},
			CreateTime: now,
			UpdateTime: now,
			CreateBy:   id,
			UpdateBy:   id,
		}
		_, err = l.svcCtx.BsUserModel.Insert(l.ctx, newBsUser)
		if err != nil {
			return nil, err
		}
		return &user.BindPhoneResp{
			UserId: id,
		}, nil
	}

	// 如果手机号已注册 将该微信绑定到手机用户 并返回手机用户ID
	err = l.svcCtx.BsUserModel.BindPhone(l.ctx, in.Phone, bsUser.Id)
	if err != nil {
		return nil, err
	}
	return &user.BindPhoneResp{
		UserId: bsUser.Id,
	}, nil
}
