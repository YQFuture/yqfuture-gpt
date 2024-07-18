package userlogic

import (
	"context"

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

// 绑定手机号
func (l *BindPhoneLogic) BindPhone(in *user.BindPhoneReq) (*user.BindPhoneResp, error) {
	// todo: add your logic here and delete this line

	return &user.BindPhoneResp{}, nil
}
