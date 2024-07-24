package user

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetMessageListLogic 获取消息列表
func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessageListLogic) GetMessageList(req *types.MessageListReq) (resp *types.MessageListResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.MessageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取失败",
			},
		}, nil
	}
	messageId, err := strconv.ParseInt(req.MessageId, 10, 64)
	if err != nil {
		l.Logger.Error("获取消息ID失败", err)
		return &types.MessageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取失败",
			},
		}, nil
	}

	// 调用RPC接口 获取消息列表
	messageListResp, err := l.svcCtx.UserClient.GetMessageList(l.ctx, &user.MessageListReq{
		UserId:     userId,
		MessageId:  messageId,
		TimeVector: req.TimeVector,
	})
	if err != nil {
		l.Logger.Error("获取消息列表失败", err)
		return &types.MessageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取失败",
			},
		}, nil
	}

	var messageList []types.MessageInfo
	for _, message := range messageListResp.Result {
		messageList = append(messageList, types.MessageInfo{
			MessageId:          strconv.FormatInt(message.MessageId, 10),
			DealFlag:           message.DealFlag,
			IgnoreFlag:         message.IgnoreFlag,
			MessageContent:     message.MessageContent,
			MessageContentType: message.MessageContentType,
			CreateTime:         message.CreateTime,
		})
	}

	return &types.MessageListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "获取成功",
		},
		Data: messageList,
	}, nil
}
