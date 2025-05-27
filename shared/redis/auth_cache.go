package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cushydigit/nanobank/shared/config"
	"github.com/cushydigit/nanobank/shared/types"
)

type AuthCacher interface {
	SetAuth(ctx context.Context, userID string, tokens *types.JWTTokens) error
	DelAuth(ctx context.Context, userID string) error
}

func authKey(userID string) string {
	return fmt.Sprintf("refresh:%s", userID)
}

func (r *MyRedisClient) SetAuth(ctx context.Context, userID string, tokens *types.JWTTokens) error {
	key := authKey(userID)
	jsonData, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	return r.rds.Set(ctx, key, jsonData, config.TTL_REFRESH_TOKEN).Err()
}

func (r *MyRedisClient) DelAuth(ctx context.Context, userID string) error {
	key := authKey(userID)
	return r.rds.Del(ctx, key).Err()
}
