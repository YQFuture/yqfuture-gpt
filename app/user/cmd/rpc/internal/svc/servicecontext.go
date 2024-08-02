package svc

import (
	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"yufuture-gpt/app/user/cmd/rpc/internal/config"
	model "yufuture-gpt/app/user/model/mongo"
	"yufuture-gpt/app/user/model/orm"
)

type ServiceContext struct {
	Config config.Config
	// Redis
	Redis *redis.Redis
	// MongoDB
	DborgpermissionModel model.DborgpermissionModel
	// MySQL模型
	BsUserModel           orm.BsUserModel
	BsOrganizationModel   orm.BsOrganizationModel
	BsUserOrgModel        orm.BsUserOrgModel
	BsMessageModel        orm.BsMessageModel
	BsMessageContentModel orm.BsMessageContentModel
	BsPermTemplateModel   orm.BsPermTemplateModel
	BsShopModel           orm.BsShopModel
	// 雪花算法
	SnowFlakeNode *snowflake.Node
}

func NewServiceContext(c config.Config) *ServiceContext {
	// MySQL
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	// 雪花算法
	snowflakeNode, err := snowflake.NewNode(c.SnowFlakeNode)
	if err != nil {
		logx.Error("雪花算法初始化失败", err)
		panic(err)
	}
	return &ServiceContext{
		Config:               c,
		Redis:                redis.MustNewRedis(c.RedisConf),
		DborgpermissionModel: model.NewDborgpermissionModel(c.Mongo.Url, c.Mongo.Database, c.Mongo.Dborgpermission),

		BsUserModel:           orm.NewBsUserModel(sqlConn),
		BsOrganizationModel:   orm.NewBsOrganizationModel(sqlConn),
		BsUserOrgModel:        orm.NewBsUserOrgModel(sqlConn),
		BsMessageModel:        orm.NewBsMessageModel(sqlConn),
		BsMessageContentModel: orm.NewBsMessageContentModel(sqlConn),
		BsPermTemplateModel:   orm.NewBsPermTemplateModel(sqlConn),
		BsShopModel:           orm.NewBsShopModel(sqlConn),

		SnowFlakeNode: snowflakeNode,
	}
}
