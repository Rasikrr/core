package redis

import (
	"context"
	"time"

	"github.com/samber/lo"
)

func (c *Cache) GetRaw(ctx context.Context, key string) *StringCMD {
	k := c.genKey(key)
	return c.client.Get(ctx, k)
}

func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Result()
}

func (c *Cache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Bytes()
}

func (c *Cache) GetBool(ctx context.Context, key string) (bool, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Bool()
}

func (c *Cache) GetInt(ctx context.Context, key string) (int, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Int()
}

func (c *Cache) GetInt64(ctx context.Context, key string) (int64, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Int64()
}

func (c *Cache) GetFloat32(ctx context.Context, key string) (float32, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Float32()
}
func (c *Cache) GetFloat64(ctx context.Context, key string) (float64, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Float64()
}

func (c *Cache) GetTime(ctx context.Context, key string) (time.Time, error) {
	k := c.genKey(key)
	return c.client.Get(ctx, k).Time()
}

func (c *Cache) MGet(ctx context.Context, keys ...string) ([]any, error) {
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
func (c *Cache) GetSet(ctx context.Context, key string, value any) (string, error) {
	k := c.genKey(key)
	return c.client.GetSet(ctx, k, value).Result()
}

// GetDel atomically gets and deletes a key
func (c *Cache) GetDel(ctx context.Context, key string) (string, error) {
	k := c.genKey(key)
	return c.client.GetDel(ctx, k).Result()
}

// GetEx gets the value and optionally sets expiration
func (c *Cache) GetEx(ctx context.Context, key string, expiration time.Duration) (string, error) {
	k := c.genKey(key)
	return c.client.GetEx(ctx, k, expiration).Result()
}
