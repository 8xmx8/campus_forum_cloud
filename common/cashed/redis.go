package cashed

import (
	"context"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(adders []string, pass, master string, db int) redis.UniversalClient {
	Client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      adders,
		Password:   pass,
		DB:         db,
		MasterName: master,
	})

	if err := redisotel.InstrumentTracing(Client); err != nil {
		panic(err)
	}

	if err := redisotel.InstrumentMetrics(Client); err != nil {
		panic(err)
	}
	if err := Client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return Client
}
