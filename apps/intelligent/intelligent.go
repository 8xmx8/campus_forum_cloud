package main

import (
	"campus_forum_cloud/apps/intelligent/internal/config"
	"campus_forum_cloud/apps/intelligent/internal/svc"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "apps/intelligent/etc/intelligent.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	svcCtx := svc.NewServiceContext(c)

	svcCtx.DAO.InitDBTable()
}
