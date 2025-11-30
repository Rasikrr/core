package redis

import (
	"context"
	"errors"
	coreCache "github.com/Rasikrr/core/cache"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

func (c *Cache) Get(ctx context.Context, key string) (any, error) {
	k := c.genKey(key)
	val, err := c.client.Get(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return val, nil
}

func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	k := c.genKey(key)
	val, err := c.client.Get(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return val, nil
}

func (c *Cache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	k := c.genKey(key)
	val, err := c.client.Get(ctx, k).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return val, nil
}

func (c *Cache) GetBool(ctx context.Context, key string) (bool, error) {
	k := c.genKey(key)
	val, err := c.client.Get(ctx, k).Bool()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return false, coreCache.ErrNotFound
		}
		return false, err
	}
	return val, nil
}

func (c *Cache) GetInt(ctx context.Context, key string) (int, error) {
	k := c.genKey(key)
	val, err := c.client.Get(ctx, k).Int()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return val, nil
}

func (c *Cache) GetInt64(ctx context.Context, key string) (int64, error) {
	k := c.genKey(key)
	val, err := c.client.Get(ctx, k).Int64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return val, nil
}

func (c *Cache) GetFloat32(ctx context.Context, key string) (float32, error) {
	k := c.genKey(key)
	val, err := c.client.Get(ctx, k).Float32()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
	}
	return val, err
}
func (c *Cache) GetFloat64(ctx context.Context, key string) (float64, error) {
	k := c.genKey(key)
	val, err := c.client.Get(ctx, k).Float64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
	}
	return val, err
}

func (c *Cache) GetTime(ctx context.Context, key string) (time.Time, error) {
	k := c.genKey(key)
	val, err := c.client.Get(ctx, k).Time()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return time.Time{}, coreCache.ErrNotFound
		}
	}
	return val, err
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
