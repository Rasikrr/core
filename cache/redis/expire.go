package redis

import (
	"context"
	"time"
)

// Expire sets a timeout on a key
func (c *Cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	k := c.genKey(key)
	return c.client.Expire(ctx, k, expiration).Err()
}

// ExpireAt sets an absolute Unix timestamp expiration on a key
func (c *Cache) ExpireAt(ctx context.Context, key string, tm time.Time) error {
	k := c.genKey(key)
	return c.client.ExpireAt(ctx, k, tm).Err()
}

// TTL returns the remaining time to live of a key
func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	k := c.genKey(key)
	return c.client.TTL(ctx, k).Result()
}

// Persist removes the expiration from a key
func (c *Cache) Persist(ctx context.Context, key string) error {
	k := c.genKey(key)
	return c.client.Persist(ctx, k).Err()
}
