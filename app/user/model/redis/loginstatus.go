package redis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const LoginStatusKey = "login:status:"

// SetLoginUser 保存登录用户
func SetLoginUser(ctx context.Context, redis *redis.Redis, userId string, seconds int) error {
	err := redis.SetexCtx(ctx, LoginStatusKey+userId, userId, seconds)
	if err != nil {
		return err
	}
	return nil
}

// GetLoginUser 获取登录用户
func GetLoginUser(ctx context.Context, redis *redis.Redis, userId string) (string, error) {
	userId, err := redis.GetCtx(ctx, LoginStatusKey+userId)
	if err != nil {
		return "", err
	}
	return userId, nil
}

// DelLoginUser 删除登录用户
func DelLoginUser(ctx context.Context, redis *redis.Redis, userId string) error {
	_, err := redis.DelCtx(ctx, LoginStatusKey+userId)
	if err != nil {
		return err
	}
	return nil
}
