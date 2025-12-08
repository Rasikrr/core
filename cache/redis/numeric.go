package redis

import "context"

func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	key = c.genKey(key)
	return c.client.Incr(ctx, key).Result()
}

func (c *Client) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	key = c.genKey(key)
	return c.client.IncrBy(ctx, key, delta).Result()
}

func (c *Client) IncrByFloat(ctx context.Context, key string, delta float64) (float64, error) {
	key = c.genKey(key)
	return c.client.IncrByFloat(ctx, key, delta).Result()
}

func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	key = c.genKey(key)
	return c.client.Decr(ctx, key).Result()
}

func (c *Client) DecrBy(ctx context.Context, key string, delta int64) (int64, error) {
	key = c.genKey(key)
	return c.client.DecrBy(ctx, key, delta).Result()
}

func (c *Client) DecrByFloat(ctx context.Context, key string, delta float64) (float64, error) {
	key = c.genKey(key)
	return c.client.IncrByFloat(ctx, key, -delta).Result()
}
