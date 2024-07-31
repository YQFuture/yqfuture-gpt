package redis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const UserOrgPrefix = "user:org:"

// AcquireDistributedLock 获取分布式锁
func AcquireDistributedLock(ctx context.Context, redis *redis.Redis, key string, seconds int) (bool, error) {
	val, err := redis.SetnxExCtx(ctx, UserOrgPrefix+key, "locked", seconds)
	if err != nil {
		return false, err
	}
	return val, nil
}

// ReleaseDistributedLock 释放分布式锁
func ReleaseDistributedLock(ctx context.Context, redis *redis.Redis, key string) error {
	_, err := redis.DelCtx(ctx, UserOrgPrefix+key)
	return err
}
