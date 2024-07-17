package loginlogic

import (
	"context"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 登录
func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 从用户表查找用户 并判断用户是否存在
	bsUser, err := l.svcCtx.BsUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		return nil, err
	}
	if bsUser == nil {
		return &user.LoginResp{
			Code: consts.PhoneIsNotRegistered,
		}, nil
	}

	// 返回用户信息
	return &user.LoginResp{
		Code: consts.Success,
		Result: &user.UserInfo{
			Id:       bsUser.Id,
			Phone:    bsUser.Phone.String,
			NickName: bsUser.NickName.String,
			HeadImg:  bsUser.HeadImg.String,
		},
	}, nil
}
