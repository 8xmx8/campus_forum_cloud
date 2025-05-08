package svc

import (
	"campus_forum_cloud/apps/intelligent/internal/config"
	"campus_forum_cloud/apps/intelligent/internal/model"
	"campus_forum_cloud/common"
	"campus_forum_cloud/common/sql"
	"github.com/importcjj/sensitive"
)

type ServiceContext struct {
	Config   config.Config
	DAO      *model.DAO
	SnFilter *sensitive.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {

	filter, err := common.InitSensitiveFilter()
	if err != nil {
		panic("init sensitive err: " + err.Error())
	}
	svc := &ServiceContext{
		Config:   c,
		DAO:      model.New(sql.Dail(c.MysqlDSN)),
		SnFilter: filter,
	}

	err = svc.DAO.InitDBTable()
	if err != nil {
		panic("init database err: " + err.Error())
	}
	return svc
}
