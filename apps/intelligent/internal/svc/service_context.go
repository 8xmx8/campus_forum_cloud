package svc

import (
	"campus_forum_cloud/apps/intelligent/internal/config"
	"campus_forum_cloud/apps/intelligent/internal/model"
	"campus_forum_cloud/common/sql"
)

type ServiceContext struct {
	Config config.Config
	DAO    *model.DAO
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DAO:    model.New(sql.Dail(c.MysqlDSN)),
	}
}
