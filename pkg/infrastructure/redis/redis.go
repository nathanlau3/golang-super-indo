package redis

import (
	"context"
	"log"

	"super-indo-api/pkg/config"

	goredis "github.com/redis/go-redis/v9"
)

func NewConnection(cfg *config.Config) *goredis.Client {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	// ping to check connection, tapi app tetap jalan meski redis mati
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("warning: redis not available: %v", err)
	}

	return rdb
}
