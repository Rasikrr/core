package redis

import "context"

func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	key = c.genKey(key)
	return c.client.Incr(ctx, key).Result()
}

func (c *Cache) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	key = c.genKey(key)
	return c.client.IncrBy(ctx, key, delta).Result()
}

func (c *Cache) IncrByFloat(ctx context.Context, key string, delta float64) (float64, error) {
	key = c.genKey(key)
	return c.client.IncrByFloat(ctx, key, delta).Result()
}

func (c *Cache) Decr(ctx context.Context, key string) (int64, error) {
	key = c.genKey(key)
	return c.client.Decr(ctx, key).Result()
}

func (c *Cache) DecrBy(ctx context.Context, key string, delta int64) (int64, error) {
	key = c.genKey(key)
	return c.client.DecrBy(ctx, key, delta).Result()
}

func (c *Cache) DecrByFloat(ctx context.Context, key string, delta float64) (float64, error) {
	key = c.genKey(key)
	return c.client.IncrByFloat(ctx, key, -delta).Result()
}
