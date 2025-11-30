package redis

import (
	"context"
	"errors"
	"time"

	coreCache "github.com/Rasikrr/core/cache"
	goredis "github.com/redis/go-redis/v9"
)

func (c *Cache) HSet(ctx context.Context, key string, values ...interface{}) error {
	key = c.genKey(key)
	return c.client.HSet(ctx, key, values...).Err()
}

func (c *Cache) HSetNX(ctx context.Context, key, field string, value interface{}) (bool, error) {
	key = c.genKey(key)
	return c.client.HSetNX(ctx, key, field, value).Result()
}

func (c *Cache) HGet(ctx context.Context, key, field string) (any, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (c *Cache) HGetBytes(ctx context.Context, key, field string) ([]byte, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (c *Cache) HGetString(ctx context.Context, key, field string) (string, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (c *Cache) HGetBool(ctx context.Context, key, field string) (bool, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Bool()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return false, coreCache.ErrNotFound
		}
		return false, err
	}
	return result, nil
}

func (c *Cache) HGetInt(ctx context.Context, key, field string) (int, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Int()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (c *Cache) HGetInt64(ctx context.Context, key, field string) (int64, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Int64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (c *Cache) HGetFloat32(ctx context.Context, key, field string) (float32, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Float32()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (c *Cache) HGetFloat64(ctx context.Context, key, field string) (float64, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Float64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (c *Cache) HGetTime(ctx context.Context, key, field string) (time.Time, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Time()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return time.Time{}, coreCache.ErrNotFound
		}
		return time.Time{}, err
	}
	return result, nil
}

func (c *Cache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	key = c.genKey(key)
	return c.client.HGetAll(ctx, key).Result()
}

func (c *Cache) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	key = c.genKey(key)
	return c.client.HMGet(ctx, key, fields...).Result()
}

func (c *Cache) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	key = c.genKey(key)
	return c.client.HDel(ctx, key, fields...).Result()
}

func (c *Cache) HExists(ctx context.Context, key, field string) (bool, error) {
	key = c.genKey(key)
	return c.client.HExists(ctx, key, field).Result()
}

func (c *Cache) HLen(ctx context.Context, key string) (int64, error) {
	key = c.genKey(key)
	return c.client.HLen(ctx, key).Result()
}

func (c *Cache) HKeys(ctx context.Context, key string) ([]string, error) {
	key = c.genKey(key)
	return c.client.HKeys(ctx, key).Result()
}

func (c *Cache) HVals(ctx context.Context, key string) ([]string, error) {
	key = c.genKey(key)
	return c.client.HVals(ctx, key).Result()
}

func (c *Cache) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	key = c.genKey(key)
	return c.client.HIncrBy(ctx, key, field, incr).Result()
}

func (c *Cache) HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error) {
	key = c.genKey(key)
	return c.client.HIncrByFloat(ctx, key, field, incr).Result()
}
