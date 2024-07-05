package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Consul consul.Conf
	// mysql
	DB struct {
		DataSource string
	}
	// mongo
	Mongo struct {
		Url                   string
		Database              string
		Kfgptaccountsentities string
	}
	// kafka生产者
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
	// kafka消费者
	KqConsumerConf kq.KqConf
	// 训练商品相关自定义配置
	TrainingGoodsConf struct {
		// 商品训练队列消费者的延迟时间，默认10秒
		ConsumeDelay int64
		// gpt图片生成文字地址
		GptImageURL string
	}
	// elasticsearch
	Elasticsearch struct {
		Addresses []string
		Username  string
		Password  string
	}
	// 雪花算法
	SnowFlakeNode int64
}
