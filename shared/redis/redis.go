package redis

import (
	"context"
	"log"

	rds "github.com/redis/go-redis/v9"
)

var Client *rds.Client

func Init(ctx context.Context, Addr string) {
	Client = rds.NewClient(&rds.Options{
		Addr:     Addr,
		Password: "",
		DB:       0,
	})

	if err := Client.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
}
