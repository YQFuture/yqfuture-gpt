package login

import (
	"context"
	"errors"
	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"
	"yufuture-gpt/app/user/model/redis"

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
	if req.Event == "subscribe" || req.Event == "SCAN" {
		if req.Ticket != "" && req.FromUserName != "" {
			err = redis.SetTicketAndOpenId(l.ctx, l.svcCtx.Redis, req.Ticket, req.FromUserName)
			if err != nil {
				return nil, err
			}
		}
	}
	return
}
