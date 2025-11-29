package redis

import (
	"context"
	"errors"
	"time"

	coreCache "github.com/Rasikrr/core/cache"
	goredis "github.com/redis/go-redis/v9"
)

func (r *cache) HSet(ctx context.Context, key string, values ...interface{}) error {
	key = r.genKey(key)
	return r.client.HSet(ctx, key, values...).Err()
}

func (r *cache) HSetNX(ctx context.Context, key, field string, value interface{}) (bool, error) {
	key = r.genKey(key)
	return r.client.HSetNX(ctx, key, field, value).Result()
}

func (r *cache) HGet(ctx context.Context, key, field string) (any, error) {
	key = r.genKey(key)
	result, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (r *cache) HGetBytes(ctx context.Context, key, field string) ([]byte, error) {
	key = r.genKey(key)
	result, err := r.client.HGet(ctx, key, field).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (r *cache) HGetString(ctx context.Context, key, field string) (string, error) {
	key = r.genKey(key)
	result, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (r *cache) HGetBool(ctx context.Context, key, field string) (bool, error) {
	key = r.genKey(key)
	result, err := r.client.HGet(ctx, key, field).Bool()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return false, coreCache.ErrNotFound
		}
		return false, err
	}
	return result, nil
}

func (r *cache) HGetInt(ctx context.Context, key, field string) (int, error) {
	key = r.genKey(key)
	result, err := r.client.HGet(ctx, key, field).Int()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (r *cache) HGetInt64(ctx context.Context, key, field string) (int64, error) {
	key = r.genKey(key)
	result, err := r.client.HGet(ctx, key, field).Int64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (r *cache) HGetFloat32(ctx context.Context, key, field string) (float32, error) {
	key = r.genKey(key)
	result, err := r.client.HGet(ctx, key, field).Float32()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (r *cache) HGetFloat64(ctx context.Context, key, field string) (float64, error) {
	key = r.genKey(key)
	result, err := r.client.HGet(ctx, key, field).Float64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (r *cache) HGetTime(ctx context.Context, key, field string) (time.Time, error) {
	key = r.genKey(key)
	result, err := r.client.HGet(ctx, key, field).Time()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return time.Time{}, coreCache.ErrNotFound
		}
		return time.Time{}, err
	}
	return result, nil
}

func (r *cache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	key = r.genKey(key)
	return r.client.HGetAll(ctx, key).Result()
}

func (r *cache) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	key = r.genKey(key)
	return r.client.HMGet(ctx, key, fields...).Result()
}

func (r *cache) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	key = r.genKey(key)
	return r.client.HDel(ctx, key, fields...).Result()
}

func (r *cache) HExists(ctx context.Context, key, field string) (bool, error) {
	key = r.genKey(key)
	return r.client.HExists(ctx, key, field).Result()
}

func (r *cache) HLen(ctx context.Context, key string) (int64, error) {
	key = r.genKey(key)
	return r.client.HLen(ctx, key).Result()
}

func (r *cache) HKeys(ctx context.Context, key string) ([]string, error) {
	key = r.genKey(key)
	return r.client.HKeys(ctx, key).Result()
}

func (r *cache) HVals(ctx context.Context, key string) ([]string, error) {
	key = r.genKey(key)
	return r.client.HVals(ctx, key).Result()
}

func (r *cache) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	key = r.genKey(key)
	return r.client.HIncrBy(ctx, key, field, incr).Result()
}

func (r *cache) HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error) {
	key = r.genKey(key)
	return r.client.HIncrByFloat(ctx, key, field, incr).Result()
}
