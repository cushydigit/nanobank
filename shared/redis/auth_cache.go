package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cushydigit/nanobank/shared/config"
	"github.com/cushydigit/nanobank/shared/types"
)

func authKey(userID string) string {
	return fmt.Sprintf("refresh:%s", userID)
}

func SetAuth(ctx context.Context, userID string, tokens *types.JWTTokens) error {
	key := authKey(userID)
	jsonData, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	return Client.Set(ctx, key, jsonData, config.TTL_REFRESH_TOKEN).Err()
}

func DelAuth(ctx context.Context, userID string) error {
	key := authKey(userID)
	return Client.Del(ctx, key).Err()
}
