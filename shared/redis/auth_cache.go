package redis

import (
	"context"
	"fmt"

	"github.com/cushydigit/nanobank/shared/config"
)

type AuthCacher interface {
	SetAuth(ctx context.Context, userID string, refreshToken string) error
	DelAuth(ctx context.Context, userID string) error
}

func authKey(userID string) string {
	return fmt.Sprintf("refresh:%s", userID)
}

func (r *MyRedisClient) SetAuth(ctx context.Context, userID string, refreshTokens string) error {
	key := authKey(userID)
	return r.rds.Set(ctx, key, refreshTokens, config.TTL_REFRESH_TOKEN).Err()
}

func (r *MyRedisClient) DelAuth(ctx context.Context, userID string) error {
	key := authKey(userID)
	return r.rds.Del(ctx, key).Err()
}
