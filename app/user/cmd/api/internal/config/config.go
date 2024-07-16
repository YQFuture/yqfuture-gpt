package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	RedisConf redis.RedisConf
	SmsConf   struct {
		AccessKeyId     string
		AccessKeySecret string
		Domain          string
		SignName        string
		TemplateCode    string
	}
}
