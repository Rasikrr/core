package redis

import (
	"context"
	"errors"
	"time"
)

func (r *cache) Set(ctx context.Context, key string, value any) error {
	k := r.genKey(key)
	return r.client.Set(ctx, k, value, 0).Err()
}

func (r *cache) MSet(ctx context.Context, keyValues ...any) error {
	if len(keyValues)%2 != 0 {
		return errors.New("invalid keyValues: must be even")
	}
	for i := 0; i < len(keyValues); i += 2 {
		keyStr, ok := keyValues[i].(string)
		if !ok {
			return errors.New("invalid key: must be string")
		}
		keyValues[i] = r.genKey(keyStr)
	}
	return r.client.MSet(ctx, keyValues...).Err()
}

func (r *cache) SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error {
	k := r.genKey(key)
	return r.client.Set(ctx, k, value, expiration).Err()
}
