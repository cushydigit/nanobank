package redis

import (
	"context"
	"fmt"

	"github.com/cushydigit/nanobank/shared/config"
)

type TokenCacher interface {
	GetToken(ctx context.Context, token string) (string, error)
	SetToken(ctx context.Context, token, txID string) error
	DelToken(ctx context.Context, token string) error
}

func tokenKey(token string) string {
	return fmt.Sprintf("confirmationToken:%s", token)
}

func (r *MyRedisClient) GetToken(ctx context.Context, token string) (string, error) {
	key := tokenKey(token)
	return r.rds.Get(ctx, key).Result()
}

func (r *MyRedisClient) SetToken(ctx context.Context, token, txID string) error {
	key := tokenKey(token)
	return r.rds.Set(ctx, key, txID, config.TTL_CONFIRMATION_TOKEN).Err()
}

func (r *MyRedisClient) DelToken(ctx context.Context, token string) error {
	key := authKey(token)
	return r.rds.Del(ctx, key).Err()
}
