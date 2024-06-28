package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"yufuture-gpt/app/training/cmd/rpc/internal/config"
	"yufuture-gpt/app/training/model/orm"
)

type ServiceContext struct {
	Config             config.Config
	TsShopModel        orm.TsShopModel
	TsGoodsModel       orm.TsGoodsModel
	TsTrainingLogModel orm.TsTrainingLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:             c,
		TsShopModel:        orm.NewTsShopModel(sqlConn),
		TsGoodsModel:       orm.NewTsGoodsModel(sqlConn),
		TsTrainingLogModel: orm.NewTsTrainingLogModel(sqlConn),
	}
}
