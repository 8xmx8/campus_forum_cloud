package svc

import (
	"campus_forum_cloud/apps/campus/internal/config"
	"campus_forum_cloud/apps/campus/internal/model"
	"campus_forum_cloud/common/cashed"
	"campus_forum_cloud/common/sql"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ServiceContext struct {
	Config config.Config
	DAO    *model.DAO
	G      *gin.Engine
}

func NewServiceContext(c config.Config) *ServiceContext {
	cc := cashed.NewRedisClient(c.Redis.Addrs, c.Redis.Pass, c.Redis.Master, c.Redis.DB)
	svc := &ServiceContext{
		DAO:    model.New(sql.Dail(c.MysqlDSN), &cc),
		G:      gin.New(),
		Config: c,
	}
	if err := svc.DAO.InitTable(); err != nil {
		panic("init table failed : " + err.Error())
	}
	return svc
}
func (s *ServiceContext) Run() error {
	s.G.GET("/ping", s.Ping)

	userGroup := s.G.Group("/user")
	{
		userGroup.POST("/register", s.Register)

	}
	err := s.G.Run(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port))
	if err != nil {
		return err
	}
	return nil
}
