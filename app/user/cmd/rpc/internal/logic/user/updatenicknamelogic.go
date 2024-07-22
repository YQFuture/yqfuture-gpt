package userlogic

import (
	"context"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNickNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateNickNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNickNameLogic {
	return &UpdateNickNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateNickName 更新昵称
func (l *UpdateNickNameLogic) UpdateNickName(in *user.UpdateNickNameReq) (*user.UpdateNickNameResp, error) {
	err := l.svcCtx.BsUserModel.UpdateNickName(l.ctx, in.NickName, in.UserId)
	if err != nil {
		l.Logger.Error("更新昵称失败", err)
		return nil, err
	}
	return &user.UpdateNickNameResp{}, nil
}
