package svc

import (
	"github.com/bwmarrin/snowflake"
	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"yufuture-gpt/app/training/cmd/rpc/internal/config"
	model "yufuture-gpt/app/training/model/mongo"
	"yufuture-gpt/app/training/model/orm"
)

type ServiceContext struct {
	Config config.Config
	// kafka生产者
	KqPusherClient *kq.Pusher
	// elasticsearch
	Elasticsearch *elastic.Client
	// 雪花算法
	SnowFlakeNode *snowflake.Node
	// mysql模型
	TsShopModel     orm.TsShopModel
	TsGoodsModel    orm.TsGoodsModel
	TsShopLogModel  orm.TsShopLogModel
	TsGoodsLogModel orm.TsGoodsLogModel
	BsDictTypeModel orm.BsDictTypeModel
	BsDictInfoModel orm.BsDictInfoModel
	// mongo模型
	ShoptrainingshoptitlesModel model.ShoptrainingshoptitlesModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化mysql连接
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	// 初始化elasticsearch连接
	esClient, err := elastic.NewClient(elastic.SetURL(c.Elasticsearch.Addresses...),
		elastic.SetBasicAuth(c.Elasticsearch.Username, c.Elasticsearch.Password),
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	// 初始化雪花算法 分布式id生成器
	snowflakeNode, err := snowflake.NewNode(c.SnowFlakeNode)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:                      c,
		KqPusherClient:              kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic, kq.WithAllowAutoTopicCreation()),
		Elasticsearch:               esClient,
		SnowFlakeNode:               snowflakeNode,
		TsShopModel:                 orm.NewTsShopModel(sqlConn),
		TsGoodsModel:                orm.NewTsGoodsModel(sqlConn),
		TsShopLogModel:              orm.NewTsShopLogModel(sqlConn),
		TsGoodsLogModel:             orm.NewTsGoodsLogModel(sqlConn),
		BsDictTypeModel:             orm.NewBsDictTypeModel(sqlConn),
		BsDictInfoModel:             orm.NewBsDictInfoModel(sqlConn),
		ShoptrainingshoptitlesModel: model.NewShoptrainingshoptitlesModel(c.Mongo.Url, c.Mongo.Database, c.Mongo.Shoptrainingshoptitles),
	}
}
