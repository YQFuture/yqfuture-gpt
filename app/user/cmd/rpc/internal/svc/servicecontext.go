package svc

import (
	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"yufuture-gpt/app/user/cmd/rpc/internal/config"
	"yufuture-gpt/app/user/model/orm"
)

type ServiceContext struct {
	Config config.Config
	// MySQL模型
	BsUserModel orm.BsUserModel
	// 雪花算法
	SnowFlakeNode *snowflake.Node
}

func NewServiceContext(c config.Config) *ServiceContext {
	// MySQL
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	snowflakeNode, err := snowflake.NewNode(c.SnowFlakeNode)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:        c,
		BsUserModel:   orm.NewBsUserModel(sqlConn),
		SnowFlakeNode: snowflakeNode,
	}
}
