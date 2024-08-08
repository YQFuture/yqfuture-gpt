package orglogic

import (
	"context"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatDetailLogic {
	return &GetChatDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取聊天记录
func (l *GetChatDetailLogic) GetChatDetail(in *user.ChatDetailReq) (*user.ChatDetailResp, error) {
	// todo: add your logic here and delete this line

	return &user.ChatDetailResp{}, nil
}
