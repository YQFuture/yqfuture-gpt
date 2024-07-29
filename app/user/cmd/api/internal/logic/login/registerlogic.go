package login

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
	"yufuture-gpt/app/user/cmd/rpc/client/login"
	"yufuture-gpt/app/user/model/redis"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRegisterLogic 注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 验证图像验证码
	answer, err := redis.GetImgCaptcha(l.ctx, l.svcCtx.Redis, req.CaptchaId)
	if err != nil {
		l.Logger.Error("从Redis中获取图像验证码答案失败", err)
		return &types.RegisterResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "注册失败",
			},
		}, nil
	}
	if answer != req.Answer {
		return &types.RegisterResp{
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
		return &types.RegisterResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "注册失败 请重试",
			},
		}, nil
	}
	if code != req.VerificationCode {
		return &types.RegisterResp{
			BaseResp: types.BaseResp{
				Code: consts.IncorrectVerificationCode,
				Msg:  "手机短信验证码错误",
			},
		}, nil
	}

	// 调用RPC接口完成注册
	registerResp, err := l.svcCtx.LoginClient.Register(l.ctx, &login.RegisterReq{
		Phone: req.Phone,
	})
	if err != nil {
		return &types.RegisterResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "注册失败 请重试",
			},
		}, nil
	}
	if registerResp.Code == consts.PhoneIsRegistered {
		return &types.RegisterResp{
			BaseResp: types.BaseResp{
				Code: consts.PhoneIsRegistered,
				Msg:  "手机号已注册 请直接登录",
			},
		}, nil
	}

	// 生成 Token
	userId := registerResp.Result.Id
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	payload := map[string]interface{}{
		"id":      userId,
		"ex_time": time.Now().AddDate(0, 0, 7),
	}
	token, err := GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, accessExpire, payload)
	if err != nil {
		l.Logger.Error("生成token失败", err)
		return &types.RegisterResp{
			BaseResp: types.BaseResp{
				Code: consts.AutomaticLoginFailure,
				Msg:  "自动登录失败 请手动登录",
			},
		}, nil
	}

	// 登录成功后，将用户信息存入 Redis
	err = redis.SetLoginUser(l.ctx, l.svcCtx.Redis, strconv.FormatInt(userId, 10), accessExpire)
	if err != nil {
		l.Logger.Error("将用户信息存入Redis失败", err)
	}

	// 返回用户信息
	return &types.RegisterResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "注册成功",
		},
		Data: types.UserInfo{
			Token:    token,
			Phone:    registerResp.Result.Phone,
			NickName: registerResp.Result.NickName,
			HeadImg:  registerResp.Result.HeadImg,
		},
	}, nil
}

// GetJwtToken 生成JWT
func GetJwtToken(secretKey string, expire int, payload map[string]interface{}) (string, error) {
	iat := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + int64(expire)
	claims["iat"] = iat
	for k, v := range payload {
		claims[k] = v
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
