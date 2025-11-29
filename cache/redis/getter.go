package redis

import (
	"context"
	"errors"
	coreCache "github.com/Rasikrr/core/cache"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

func (r *cache) Get(ctx context.Context, key string) (any, error) {
	k := r.genKey(key)
	val, err := r.client.Get(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return val, nil
}

func (r *cache) GetString(ctx context.Context, key string) (string, error) {
	k := r.genKey(key)
	val, err := r.client.Get(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return val, nil
}

func (r *cache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	k := r.genKey(key)
	val, err := r.client.Get(ctx, k).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return val, nil
}

func (r *cache) GetBool(ctx context.Context, key string) (bool, error) {
	k := r.genKey(key)
	val, err := r.client.Get(ctx, k).Bool()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return false, coreCache.ErrNotFound
		}
		return false, err
	}
	return val, nil
}

func (r *cache) GetInt(ctx context.Context, key string) (int, error) {
	k := r.genKey(key)
	val, err := r.client.Get(ctx, k).Int()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return val, nil
}

func (r *cache) GetInt64(ctx context.Context, key string) (int64, error) {
	k := r.genKey(key)
	val, err := r.client.Get(ctx, k).Int64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return val, nil
}

func (r *cache) GetFloat32(ctx context.Context, key string) (float32, error) {
	k := r.genKey(key)
	val, err := r.client.Get(ctx, k).Float32()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
	}
	return val, err
}
func (r *cache) GetFloat64(ctx context.Context, key string) (float64, error) {
	k := r.genKey(key)
	val, err := r.client.Get(ctx, k).Float64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
	}
	return val, err
}

func (r *cache) GetTime(ctx context.Context, key string) (time.Time, error) {
	k := r.genKey(key)
	val, err := r.client.Get(ctx, k).Time()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return time.Time{}, coreCache.ErrNotFound
		}
	}
	return val, err
}

func (r *cache) MGet(ctx context.Context, keys ...string) ([]any, error) {
	keys = lo.Map(keys, func(k string, _ int) string {
		return r.genKey(k)
	})
	values, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	return values, nil
}
