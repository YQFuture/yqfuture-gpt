package svc

import (
	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"yufuture-gpt/app/user/cmd/rpc/internal/config"
	"yufuture-gpt/app/user/model/orm"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
	// MySQL模型
	BsUserModel           orm.BsUserModel
	BsOrganizationModel   orm.BsOrganizationModel
	BsUserOrgModel        orm.BsUserOrgModel
	BsMessageModel        orm.BsMessageModel
	BsMessageContentModel orm.BsMessageContentModel
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
		Config:                c,
		Redis:                 redis.MustNewRedis(c.RedisConf),
		BsUserModel:           orm.NewBsUserModel(sqlConn),
		BsOrganizationModel:   orm.NewBsOrganizationModel(sqlConn),
		BsUserOrgModel:        orm.NewBsUserOrgModel(sqlConn),
		BsMessageModel:        orm.NewBsMessageModel(sqlConn),
		BsMessageContentModel: orm.NewBsMessageContentModel(sqlConn),
		SnowFlakeNode:         snowflakeNode,
	}
}
