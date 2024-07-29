package user

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"yufuture-gpt/app/user/cmd/api/internal/logic/login"
	"yufuture-gpt/app/user/cmd/rpc/client/user"
	"yufuture-gpt/app/user/model/redis"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewBindPhoneLogic 绑定手机号码
func NewBindPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindPhoneLogic {
	return &BindPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindPhoneLogic) BindPhone(req *types.BindPhoneReq) (resp *types.BindPhoneResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.BindPhoneResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "绑定失败 请重试",
			},
		}, nil
	}
	userIdString := strconv.FormatInt(userId, 10)
	// 根据微信临时用户ID获取对应的OpenID
	openId, err := redis.GetOpenIdByTempUserId(l.ctx, l.svcCtx.Redis, userIdString)
	if err != nil {
		l.Logger.Error("从Redis中获取临时用户ID对应的OpenID失败", err)
		return &types.BindPhoneResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "绑定失败 请重试",
			},
		}, nil
	}
	err = redis.DelOpenIdByTempUserId(l.ctx, l.svcCtx.Redis, userIdString)
	if err != nil {
		l.Logger.Error("从Redis中删除临时用户ID对应的OpenID失败", err)
	}

	// 验证手机短信验证码
	code, err := redis.GetVerificationCode(l.ctx, l.svcCtx.Redis, req.Phone)
	if err != nil {
		l.Logger.Error("从Redis中获取手机短信验证码失败", err)
		return &types.BindPhoneResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "绑定失败 请重试",
			},
		}, nil
	}
	if code != req.VerificationCode {
		return &types.BindPhoneResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "短信验证码错误",
			},
		}, nil
	}

	// 调用RPC接口 绑定手机号码
	bindPhoneResp, err := l.svcCtx.UserClient.BindPhone(l.ctx, &user.BindPhoneReq{
		Openid: openId,
		Phone:  req.Phone,
	})
	if err != nil {
		return &types.BindPhoneResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "绑定失败 请重试",
			},
		}, nil
	}

	// 生成新Token
	userId = bindPhoneResp.UserId
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	if req.ThirtyDaysFreeLogin {
		accessExpire = 2592000
	}
	payload := map[string]interface{}{
		"id":      userId,
		"ex_time": time.Now().AddDate(0, 0, 7),
	}
	token, err := login.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, accessExpire, payload)
	if err != nil {
		l.Logger.Error("生成token失败", err)
		return &types.BindPhoneResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "绑定失败 请重试",
			},
		}, nil
	}

	// 登录成功后，将用户信息存入 Redis
	err = redis.SetLoginUser(l.ctx, l.svcCtx.Redis, strconv.FormatInt(userId, 10), accessExpire)

	return &types.BindPhoneResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "绑定成功",
		},
		Data: types.UserInfo{
			Token: token,
		},
	}, nil
}
