package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	// Consul
	Consul consul.Conf
	// Redis
	RedisConf redis.RedisConf
	// MongoDB
	Mongo struct {
		Url             string
		Database        string
		Dborgpermission string
		Dbuseroperation string
	}
	// MySQL
	DB struct {
		DataSource string
	}
	// 雪花算法节点ID
	SnowFlakeNode int64
}
