package userlogic

import (
	"context"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetMessageList 获取消息列表
func (l *GetMessageListLogic) GetMessageList(in *user.MessageListReq) (*user.MessageListResp, error) {
	// 获取用户信息
	bsUser, err := l.svcCtx.BsUserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	// 获取消息列表
	messageList, err := l.svcCtx.BsMessageModel.FindMessageList(l.ctx, in.UserId, bsUser.NowOrgId, in.MessageId, in.TimeVector)
	if err != nil {
		return nil, err
	}

	var messageInfoList []*user.MessageInfo
	for _, message := range *messageList {
		messageInfoList = append(messageInfoList, &user.MessageInfo{
			MessageId:          message.Id,
			MessageContentType: message.MessageContentType,
			MessageContent:     message.MessageContent,
			DealFlag:           message.DealFlag,
			IgnoreFlag:         message.IgnoreFlag,
			CreateTime:         message.CreateTime.Unix(),
		})
	}
	return &user.MessageListResp{
		Result: messageInfoList,
	}, nil
}
