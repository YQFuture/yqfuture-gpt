package redis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const LoginQrCodePrefix = "login:qrcode:ticket:"

// SetTicketAndOpenId 保存登录二维码票据和对应的OpenID到Redis
func SetTicketAndOpenId(ctx context.Context, redis *redis.Redis, ticket, openid string) error {
	err := redis.SetexCtx(ctx, LoginQrCodePrefix+ticket, openid, 120)
	if err != nil {
		return err
	}
	return nil
}

// GetOpenId 根据票据获取对应的OpenID
func GetOpenId(ctx context.Context, redis *redis.Redis, ticket string) (string, error) {
	openid, err := redis.GetCtx(ctx, LoginQrCodePrefix+ticket)
	if err != nil {
		return "", err
	}
	return openid, nil
}
