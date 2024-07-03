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
	Config                     config.Config
	TsShopModel                orm.TsShopModel
	TsGoodsModel               orm.TsGoodsModel
	TsTrainingLogModel         orm.TsTrainingLogModel
	BsDictTypeModel            orm.BsDictTypeModel
	BsDictInfoModel            orm.BsDictInfoModel
	KfgptaccountsentitiesModel model.KfgptaccountsentitiesModel
	KqPusherClient             *kq.Pusher
	Elasticsearch              *elastic.Client
	SnowFlakeNode              *snowflake.Node
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	esClient, err := elastic.NewClient(elastic.SetURL(c.Elasticsearch.Addresses...),
		elastic.SetBasicAuth(c.Elasticsearch.Username, c.Elasticsearch.Password),
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	node, err := snowflake.NewNode(c.SnowFlakeNode)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:                     c,
		TsShopModel:                orm.NewTsShopModel(sqlConn),
		TsGoodsModel:               orm.NewTsGoodsModel(sqlConn),
		TsTrainingLogModel:         orm.NewTsTrainingLogModel(sqlConn),
		BsDictTypeModel:            orm.NewBsDictTypeModel(sqlConn),
		BsDictInfoModel:            orm.NewBsDictInfoModel(sqlConn),
		KfgptaccountsentitiesModel: model.NewKfgptaccountsentitiesModel(c.Mongo.Url, c.Mongo.Database, c.Mongo.Kfgptaccountsentities),
		KqPusherClient:             kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic, kq.WithAllowAutoTopicCreation()),
		Elasticsearch:              esClient,
		SnowFlakeNode:              node,
	}
}
