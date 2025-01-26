package main

import (
	"campus_forum_cloud/apps/campus/internal/config"
	"campus_forum_cloud/apps/campus/internal/svc"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "apps/campus/etc/campus.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	logx.MustSetup(logx.LogConf{
		ServiceName: c.Name,
		Path:        c.LogPath,
		Mode:        c.LogMode,
		Encoding:    c.LogEncoding,
	})
	defer func() {
		err := logx.Close()
		if err != nil {
			panic(err)
		}
	}()
	server := svc.NewServiceContext(c)
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
