package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Consul consul.Conf
	DB     struct {
		DataSource string
	}
	Mongo struct {
		Url                   string
		Database              string
		Kfgptaccountsentities string
	}
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
	KqConsumerConf    kq.KqConf
	TrainingGoodsConf struct {
		ConsumeDelay int64
	}
	Elasticsearch struct {
		Addresses []string
		Username  string
		Password  string
	}
}
