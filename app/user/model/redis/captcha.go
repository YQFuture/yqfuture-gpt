package redis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const ImgCaptchaPrefix = "captcha:img:"

// SetImgCaptcha 保存图像验证码答案到Redis
func SetImgCaptcha(ctx context.Context, redis *redis.Redis, id, answer string) error {
	err := redis.SetexCtx(ctx, ImgCaptchaPrefix+id, answer, 65)
	if err != nil {
		return err
	}
	return nil
}

// GetImgCaptcha 根据id从Redis中获取图像验证码答案
func GetImgCaptcha(ctx context.Context, redis *redis.Redis, id string) (string, error) {
	answer, err := redis.GetCtx(ctx, ImgCaptchaPrefix+id)
	if err != nil {
		return "", err
	}
	return answer, nil
}
