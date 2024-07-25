package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UserClientConf zrpc.RpcClientConf
	// JWT
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	// Redis
	RedisConf redis.RedisConf
	// 阿里云短信服务
	SmsConf struct {
		AccessKeyId     string
		AccessKeySecret string
		Domain          string
		SignName        string
		TemplateCode    string
	}
	// 阿里云OSS服务
	OssConf struct {
		Endpoint        string
		AccessKeyId     string
		AccessKeySecret string
		BucketName      string
	}
	// 微信公众号
	WechatConf struct {
		AppId          string
		Secret         string
		AccessTokenUrl string
		TicketUrl      string
		QrCodeUrl      string
		Token          string
	}
}
