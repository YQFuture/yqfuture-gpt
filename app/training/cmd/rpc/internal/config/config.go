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
		Url                      string
		Database                 string
		Dbsavegoodscrawlertitles string
		Dbpresettingshoptitles   string
		Dbpresettinggoodstitles  string
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
		// 申请获取商品JSON接口
		ApplyGoodsJsonUrl string
		// 申请获取商品JSON接口的渠道ID
		ApplyGoodsJsonChannel string
		// 拉取商品JSON接口
		FetchGoodsJsonUrl string

		// 大模型接口-申请爬取商品ID列表
		ApplyGoodsIdListUrl string
		// 大模型接口-获取预估结果
		FetchEstimateResultUrl string
		// 大模型接口-创建批处理任务
		CreateBatchTaskUrl string
		// 大模型接口-查询处理任务状态
		QueryBatchTaskStatusUrl string
		// 大模型接口-获取批处理返回结果
		QueryBatchTaskResultUrl string
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
