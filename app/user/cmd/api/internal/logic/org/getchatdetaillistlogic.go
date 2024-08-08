package org

import (
	"context"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatDetailListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetChatDetailListLogic 获取聊天记录列表
func NewGetChatDetailListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatDetailListLogic {
	return &GetChatDetailListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatDetailListLogic) GetChatDetailList(req *types.ChatDetailReq) (resp *types.ChatDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
