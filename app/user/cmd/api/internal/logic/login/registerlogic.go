package login

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"time"
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
				Msg:  "注册失败",
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

	// TODO 发送注册RPC调用

	// TODO 判断用户是否已注册

	// 生成 Token
	payload := map[string]interface{}{
		"id":      1811692312716120064,
		"ex_time": time.Now().AddDate(0, 0, 7),
	}
	token, err := GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire, payload)
	if err != nil {
		l.Logger.Error("生成token失败", err)
		return &types.RegisterResp{
			BaseResp: types.BaseResp{
				Code: consts.AutomaticLoginFailure,
				Msg:  "自动登录失败 请手动登录",
			},
		}, nil
	}
	return &types.RegisterResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "注册成功",
		},
		Data: types.UserInfo{
			Token: token,
		},
	}, nil
}

// GetJwtToken 生成JWT
func GetJwtToken(secretKey string, expire int64, payload map[string]interface{}) (string, error) {
	iat := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + expire
	claims["iat"] = iat
	for k, v := range payload {
		claims[k] = v
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
