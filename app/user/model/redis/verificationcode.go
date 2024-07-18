package redis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const VerificationCodePrefix = "verification:code:phone:"

// SetVerificationCode 保存手机短信验证码到Redis
func SetVerificationCode(ctx context.Context, redis *redis.Redis, phone, code string) error {
	err := redis.SetexCtx(ctx, VerificationCodePrefix+phone, code, 60)
	if err != nil {
		return err
	}
	return nil
}

// GetVerificationCode 根据手机号码从Redis中获取手机短信验证码
func GetVerificationCode(ctx context.Context, redis *redis.Redis, phone string) (string, error) {
	code, err := redis.GetCtx(ctx, VerificationCodePrefix+phone)
	if err != nil {
		return "", err
	}
	return code, nil
}
