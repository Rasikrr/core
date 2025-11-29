package redis

import (
	"context"
	"time"
)

func (r *cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	k := r.genKey(key)
	return r.client.Expire(ctx, k, expiration).Err()
}
