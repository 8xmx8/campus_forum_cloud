package main

import (
	"campus_forum_cloud/apps/intelligent/internal/chatredischannel"
	"campus_forum_cloud/apps/intelligent/internal/config"
	"campus_forum_cloud/apps/intelligent/internal/svc"
	"campus_forum_cloud/common/redissubscriber"
	"context"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"os"
	"os/signal"
	"syscall"
)

var configFile = flag.String("f", "apps/intelligent/etc/intelligent.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	channels := []string{
		c.CommentChain,
		c.ArticleChain,
	}
	// 创建根上下文和取消函数
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svcCtx := svc.NewServiceContext(c)
	subscriber := redissubscriber.NewRedisSubscriber(
		c.Redis,
		channels,
		3,
		chatredischannel.NewChatRedisChannelHandler(ctx, svcCtx).Hander)
	subscriber.Start(ctx)
	defer subscriber.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞主函数，等待信号
	<-sigChan
	cancel()
}
