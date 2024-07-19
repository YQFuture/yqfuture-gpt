package login

import (
	"context"
	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/thirdparty"
	"yufuture-gpt/app/user/cmd/api/internal/types"
	"yufuture-gpt/common/consts"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoginQrCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetLoginQrCodeLogic 获取微信登录二维码
func NewGetLoginQrCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginQrCodeLogic {
	return &GetLoginQrCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoginQrCodeLogic) GetLoginQrCode(req *types.BaseReq) (resp *types.LoginQrCodeResp, err error) {
	ticketQrCodeUrl, ticket, err := thirdparty.GetWechatLoginQrCode(l.ctx, l.svcCtx.Redis,
		l.svcCtx.Config.WechatConf.AccessTokenUrl,
		l.svcCtx.Config.WechatConf.AppId,
		l.svcCtx.Config.WechatConf.Secret,
		l.svcCtx.Config.WechatConf.TicketUrl,
		l.svcCtx.Config.WechatConf.QrCodeUrl)
	if err != nil {
		l.Logger.Error("获取微信登录二维码失败: ", err)
		return &types.LoginQrCodeResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取微信登录二维码失败",
			},
		}, nil
	}
	return &types.LoginQrCodeResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "获取微信登录二维码成功",
		},
		Data: types.LoginQrCodeData{
			TicketQrCodeUrl: ticketQrCodeUrl,
			Ticket:          ticket,
		},
	}, nil
}
