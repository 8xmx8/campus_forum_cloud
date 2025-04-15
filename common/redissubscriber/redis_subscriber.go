package redissubscriber

import (
	"campus_forum_cloud/common"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
)

type RedisSubscriber struct {
	client   redis.UniversalClient    // Redis 客户端
	pubsub   *redis.PubSub            // PubSub 对象
	channels []string                 // 订阅的频道列表
	handler  func(msg *redis.Message) // 消息处理函数

	worker int // 工作协程数量
}

func NewRedisSubscriber(config *common.RedisClientOpt, channels []string, workerNum int, handler func(msg *redis.Message)) *RedisSubscriber {
	var RS redis.UniversalClient
	if config.Type == "node" {
		addrs := strings.Split(config.Addr, ",")
		RS = redis.NewClient(&redis.Options{
			Addr:     addrs[0],        // Redis 地址
			Password: config.Password, // 密码
			DB:       config.DB,       // 数据库
		})
	} else if config.Type == "cluster" {
		addrs := strings.Split(config.Addr, ",")
		RS = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,           // Redis 地址
			Password: config.Password, // 密码
		})
	} else if config.Type == "sentinel" {
		addrs := strings.Split(config.Addr, ",")
		RS = redis.NewFailoverClusterClient(&redis.FailoverOptions{
			SentinelAddrs: addrs,           // Redis 地址,哨兵集群地址
			Password:      config.Password, // 密码
		})
	}

	ctx := context.Background()
	pong, err := RS.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Redis:", pong)

	return &RedisSubscriber{
		client:   RS,
		channels: channels,
		handler:  handler,
		worker:   workerNum,
	}
}

// Start 启动订阅者
func (s *RedisSubscriber) Start(ctx context.Context) {
	// 订阅频道
	s.pubsub = s.client.Subscribe(ctx, s.channels...)

	// 获取消息通道
	msgChan := s.pubsub.Channel()

	for i := 0; i < s.worker; i++ { // 5 个并发 worker
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done(): // 监听上下文取消信号
					log.Printf("Worker %d: context cancelled, stopping...\n", workerID)
					return
				case msg, ok := <-msgChan:
					if !ok {
						log.Printf("Worker %d: message channel closed, stopping...\n", workerID)
						return
					}
					// 调用消息处理函数
					s.handler(msg)
				}
			}
		}(i)
	}

	log.Printf("Subscribed to channels: %v\n", s.channels)
}

func (s *RedisSubscriber) Stop() {
	s.pubsub.Close()
}
