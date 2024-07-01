package svc

import (
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
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:                     c,
		TsShopModel:                orm.NewTsShopModel(sqlConn),
		TsGoodsModel:               orm.NewTsGoodsModel(sqlConn),
		TsTrainingLogModel:         orm.NewTsTrainingLogModel(sqlConn),
		BsDictTypeModel:            orm.NewBsDictTypeModel(sqlConn),
		BsDictInfoModel:            orm.NewBsDictInfoModel(sqlConn),
		KfgptaccountsentitiesModel: model.NewKfgptaccountsentitiesModel(c.Mongo.Url, c.Mongo.Database, c.Mongo.Kfgptaccountsentities),
	}
}
