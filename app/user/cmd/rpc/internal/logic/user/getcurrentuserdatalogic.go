package userlogic

import (
	"context"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentUserDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCurrentUserDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserDataLogic {
	return &GetCurrentUserDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCurrentUserData 获取当前登录用户数据
func (l *GetCurrentUserDataLogic) GetCurrentUserData(in *user.CurrentUserDataReq) (*user.CurrentUserDataResp, error) {
	bsUser, err := l.svcCtx.BsUserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	// 判断是否未绑定手机号
	if bsUser.Phone.Valid == false || bsUser.Phone.String == "" {
		return &user.CurrentUserDataResp{
			Code: consts.PhoneTsNotBound,
		}, nil
	}

	return &user.CurrentUserDataResp{}, nil
}
