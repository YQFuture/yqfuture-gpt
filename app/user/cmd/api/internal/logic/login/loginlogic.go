package login

import (
	"context"
	"strconv"
	"time"
	"yufuture-gpt/app/user/cmd/rpc/client/login"
	"yufuture-gpt/app/user/model/redis"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLoginLogic 登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 验证图像验证码
	answer, err := redis.GetImgCaptcha(l.ctx, l.svcCtx.Redis, req.CaptchaId)
	if err != nil {
		l.Logger.Error("从Redis中获取图像验证码答案失败", err)
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "注册失败",
			},
		}, nil
	}
	if answer != req.Answer {
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: consts.IncorrectCaptcha,
				Msg:  "图像验证码错误",
			},
		}, nil
	}

	// 验证手机短信验证码
	code, err := redis.GetVerificationCode(l.ctx, l.svcCtx.Redis, req.Phone)
	if err != nil {
		l.Logger.Error("从Redis中获取手机短信验证码失败", err)
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "注册失败 请重试",
			},
		}, nil
	}
	if code != req.VerificationCode {
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: consts.IncorrectVerificationCode,
				Msg:  "手机短信验证码错误",
			},
		}, nil
	}

	// 调用RPC接口完成登录
	loginResp, err := l.svcCtx.LoginClient.Login(l.ctx, &login.LoginReq{
		Phone: req.Phone,
	})
	if err != nil {
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "登录失败 请重试",
			},
		}, nil
	}
	if loginResp.Code == consts.PhoneIsNotRegistered {
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: consts.PhoneIsNotRegistered,
				Msg:  "手机号未注册 请先进行注册",
			},
		}, nil
	}

	// 生成 Token
	userId := loginResp.Result.Id
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	if req.ThirtyDaysFreeLogin {
		accessExpire = 2592000
	}
	payload := map[string]interface{}{
		"id":      userId,
		"ex_time": time.Now().AddDate(0, 0, 7),
	}
	token, err := GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, accessExpire, payload)
	if err != nil {
		l.Logger.Error("生成token失败", err)
		return &types.LoginResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "登录失败 请重试",
			},
		}, nil
	}

	// 登录成功后，将用户信息存入 Redis
	err = redis.SetLoginUser(l.ctx, l.svcCtx.Redis, strconv.FormatInt(userId, 10), accessExpire)

	return &types.LoginResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "登录成功",
		},
		Data: types.UserInfo{
			Token:    token,
			Phone:    loginResp.Result.Phone,
			NickName: loginResp.Result.NickName,
			HeadImg:  loginResp.Result.HeadImg,
		},
	}, nil
}
