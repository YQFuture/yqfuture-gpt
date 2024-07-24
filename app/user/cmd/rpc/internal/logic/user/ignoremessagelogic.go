package userlogic

import (
	"context"
	"errors"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type IgnoreMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIgnoreMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IgnoreMessageLogic {
	return &IgnoreMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// IgnoreMessage 忽略消息
func (l *IgnoreMessageLogic) IgnoreMessage(in *user.IgnoreMessageReq) (*user.IgnoreMessageResp, error) {
	bsMessage, err := l.svcCtx.BsMessageModel.FindOne(l.ctx, in.MessageId)
	if err != nil {
		l.Logger.Error("根据消息ID获取消息失败", err)
		return nil, err
	}
	if bsMessage.UserId != in.UserId {
		l.Logger.Error("用户ID和消息记录对应的用户ID不一致")
		return nil, errors.New("用户ID和消息记录对应的用户ID不一致")
	}

	err = l.svcCtx.BsMessageModel.IgnoreMessage(l.ctx, in.MessageId)
	if err != nil {
		return nil, err
	}
	return &user.IgnoreMessageResp{}, nil
}
