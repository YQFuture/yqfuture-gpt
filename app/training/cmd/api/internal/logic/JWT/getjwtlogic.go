package JWT

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJWTLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetJWTLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJWTLogic {
	return &GetJWTLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJWTLogic) GetJWT(req *types.BaseReq) (resp *types.BaseResp, err error) {
	playLoad := map[string]interface{}{
		"id":      1,
		"ex_time": time.Now().AddDate(0, 0, 7),
	}
	jwtToken, err := getJwtToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, playLoad)
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code: consts.Success,
		Msg:  jwtToken,
	}, nil
}

// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func getJwtToken(secretKey string, iat, seconds int64, payload map[string]interface{}) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	for k, v := range payload {
		claims[k] = v
	}
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
