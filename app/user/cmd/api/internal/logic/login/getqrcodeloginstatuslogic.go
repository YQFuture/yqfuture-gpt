package login

import (
	"context"
	"time"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"
	"yufuture-gpt/app/user/model/redis"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQrCodeLoginStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetQrCodeLoginStatusLogic 获取微信扫码登录状态
func NewGetQrCodeLoginStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQrCodeLoginStatusLogic {
	return &GetQrCodeLoginStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQrCodeLoginStatusLogic) GetQrCodeLoginStatus(req *types.QrCodeLoginStatusReq) (resp *types.QrCodeLoginStatusResp, err error) {
	openId, err := redis.GetOpenId(l.ctx, l.svcCtx.Redis, req.Ticket)
	if err != nil {
		l.Logger.Error("从Redis中获取OpenID失败", err)
		return &types.QrCodeLoginStatusResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "登录失败",
			},
		}, nil
	}
	if openId == "" {
		l.Logger.Error("从Redis中未获取OpenID", err)
		return &types.QrCodeLoginStatusResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "登录失败",
			},
		}, nil
	}

	// 调用RPC接口 获取用户信息
	infoResp, err := l.svcCtx.LoginClient.GetWechatUserInfo(l.ctx, &user.WechatUserInfoReq{
		Openid: openId,
	})
	if err != nil {
		return &types.QrCodeLoginStatusResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "登录失败",
			},
		}, nil
	}

	// 生成 Token
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	if req.ThirtyDaysFreeLogin {
		accessExpire = 2592000
	}
	payload := map[string]interface{}{
		"id":      infoResp.Result.Id,
		"ex_time": time.Now().AddDate(0, 0, 7),
	}
	token, err := GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, accessExpire, payload)
	if err != nil {
		l.Logger.Error("生成token失败", err)
		return &types.QrCodeLoginStatusResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "登录失败 请重试",
			},
		}, nil
	}

	return &types.QrCodeLoginStatusResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "登录成功",
		},
		Data: types.UserInfo{
			Token:    token,
			Phone:    infoResp.Result.Phone,
			NickName: infoResp.Result.NickName,
			HeadImg:  infoResp.Result.HeadImg,
		},
	}, nil
}
