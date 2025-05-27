package redis

import (
	"context"
	"log"
	"sync"

	rds "github.com/redis/go-redis/v9"
)

var (
	client *MyRedisClient
	once   sync.Once
)

type MyRedisClient struct {
	rds *rds.Client
}

func MyRedisClientInit(ctx context.Context, Addr string) *MyRedisClient {
	once.Do(func() {
		c := rds.NewClient(&rds.Options{
			Addr:     Addr,
			Password: "",
			DB:       0,
		})
		client = &MyRedisClient{rds: c}
	})

	// check the connection
	if err := client.rds.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return client
}

func GetInstance() *MyRedisClient {
	return client
}
