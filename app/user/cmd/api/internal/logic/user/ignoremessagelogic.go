package user

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/user/cmd/rpc/client/user"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IgnoreMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewIgnoreMessageLogic 忽略消息
func NewIgnoreMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IgnoreMessageLogic {
	return &IgnoreMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IgnoreMessageLogic) IgnoreMessage(req *types.IgnoreMessageReq) (resp *types.BaseResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "操作失败",
		}, nil
	}
	messageId, err := strconv.ParseInt(req.MessageId, 10, 64)
	if err != nil {
		l.Logger.Error("获取消息ID失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "操作失败",
		}, nil
	}

	// 调用RPC接口 将消息置为已忽略状态
	_, err = l.svcCtx.UserClient.IgnoreMessage(l.ctx, &user.IgnoreMessageReq{
		UserId:    userId,
		MessageId: messageId,
	})
	if err != nil {
		l.Logger.Error("忽略消息失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "操作失败",
		}, nil
	}

	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "操作成功",
	}, nil
}
