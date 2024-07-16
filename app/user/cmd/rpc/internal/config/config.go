package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	// Consul
	Consul consul.Conf
	// MySQL
	DB struct {
		DataSource string
	}
	// 雪花算法节点ID
	SnowFlakeNode int64
}
