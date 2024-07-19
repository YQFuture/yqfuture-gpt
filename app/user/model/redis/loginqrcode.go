package redis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const WechatAccessTokenKey = "login:wechat:access_token"
const LoginQrCodePrefix = "login:qrcode:ticket:"
const TempUserIdPrefix = "login:temp:userid:"

// SetAccessToken 保存微信的AccessToken到Redis
func SetAccessToken(ctx context.Context, redis *redis.Redis, accessToken string) error {
	err := redis.SetexCtx(ctx, WechatAccessTokenKey, accessToken, 7000)
	if err != nil {
		return err
	}
	return nil
}

// GetAccessToken 从Redis中获取微信的AccessToken
func GetAccessToken(ctx context.Context, redis *redis.Redis) (string, error) {
	accessToken, err := redis.GetCtx(ctx, WechatAccessTokenKey)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// SetTicketAndOpenId 保存登录二维码票据和对应的OpenID到Redis
func SetTicketAndOpenId(ctx context.Context, redis *redis.Redis, ticket, openid string) error {
	err := redis.SetexCtx(ctx, LoginQrCodePrefix+ticket, openid, 120)
	if err != nil {
		return err
	}
	return nil
}

// GetOpenIdByTicket 从Redis中获取票据对应的OpenID
func GetOpenIdByTicket(ctx context.Context, redis *redis.Redis, ticket string) (string, error) {
	openid, err := redis.GetCtx(ctx, LoginQrCodePrefix+ticket)
	if err != nil {
		return "", err
	}
	return openid, nil
}

// SetTempUserIdAndOpenId 保存微信临时用户ID和对应的OpenID到Redis
func SetTempUserIdAndOpenId(ctx context.Context, redis *redis.Redis, userId, openid string) error {
	err := redis.SetexCtx(ctx, TempUserIdPrefix+userId, openid, 3600)
	if err != nil {
		return err
	}
	return nil
}

// GetOpenIdByTempUserId 从Redis中获取微信临时用户ID对应的OpenID
func GetOpenIdByTempUserId(ctx context.Context, redis *redis.Redis, userId string) (string, error) {
	openid, err := redis.GetCtx(ctx, TempUserIdPrefix+userId)
	if err != nil {
		return "", err
	}
	return openid, nil
}

// DelOpenIdByTempUserId 从Redis中删除微信临时用户ID对应的OpenID
func DelOpenIdByTempUserId(ctx context.Context, redis *redis.Redis, userId string) error {
	_, err := redis.DelCtx(ctx, TempUserIdPrefix+userId)
	if err != nil {
		return err
	}
	return nil
}
