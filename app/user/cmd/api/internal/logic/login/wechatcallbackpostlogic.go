package login

import (
	"context"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WechatCallBackPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewWechatCallBackPostLogic 微信回调
func NewWechatCallBackPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatCallBackPostLogic {
	return &WechatCallBackPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatCallBackPostLogic) WechatCallBackPost(req *types.WechatCallBackPostReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
