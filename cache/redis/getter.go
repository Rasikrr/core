package redis

import (
	"context"
	"time"

	"github.com/samber/lo"
)

func (c *Client) GetRaw(ctx context.Context, key string) *StringCMD {
	k := c.genKey(key)
	return c.client.Get(ctx, k)
}

func (c *Client) GetString(ctx context.Context, key string) (string, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Result()
}

func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Bytes()
}

func (c *Client) GetBool(ctx context.Context, key string) (bool, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Bool()
}

func (c *Client) GetInt(ctx context.Context, key string) (int, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Int()
}

func (c *Client) GetInt64(ctx context.Context, key string) (int64, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Int64()
}

func (c *Client) GetFloat32(ctx context.Context, key string) (float32, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Float32()
}
func (c *Client) GetFloat64(ctx context.Context, key string) (float64, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Float64()
}

func (c *Client) GetTime(ctx context.Context, key string) (time.Time, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Time()
}

func (c *Client) MGet(ctx context.Context, keys ...string) ([]any, error) {
	keys = lo.Map(keys, func(k string, _ int) string {
		return c.genKey(k)
	})
	values, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	return values, nil
}

// GetSet atomically sets key to value and returns the old value
func (c *Client) GetSet(ctx context.Context, key string, value any) (string, error) {
	k := c.genKey(key)
	return c.client.GetSet(ctx, k, value).Result()
}

// GetDel atomically gets and deletes a key
func (c *Client) GetDel(ctx context.Context, key string) (string, error) {
	k := c.genKey(key)
	return c.client.GetDel(ctx, k).Result()
}

// GetEx gets the value and optionally sets expiration
func (c *Client) GetEx(ctx context.Context, key string, expiration time.Duration) (string, error) {
	k := c.genKey(key)
	return c.client.GetEx(ctx, k, expiration).Result()
}
