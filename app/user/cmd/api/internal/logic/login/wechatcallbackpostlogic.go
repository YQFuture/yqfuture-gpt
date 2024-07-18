package login

import (
	"context"
	"errors"
	"fmt"

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
	if !CheckSignature(l.svcCtx.Config.WechatConf.Token, req.Signature, req.Timestamp, req.Nonce) {
		return nil, errors.New("签名错误")
	}

	fmt.Println("-------------------------------------------")
	fmt.Println(req.Event)
	fmt.Println(req.Ticket)
	fmt.Println(req.FromUserName)
	fmt.Println("-------------------------------------------")
	return
}
