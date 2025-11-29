package redis

import (
	"context"
	"errors"
	"time"

	coreCache "github.com/Rasikrr/core/cache"
	goredis "github.com/redis/go-redis/v9"
)

func (r *cache) RPush(ctx context.Context, key string, value ...any) error {
	k := r.genKey(key)
	return r.client.RPush(ctx, k, value...).Err()
}

func (r *cache) LPush(ctx context.Context, key string, value ...any) error {
	k := r.genKey(key)
	return r.client.LPush(ctx, k, value...).Err()
}

func (r *cache) LLen(ctx context.Context, key string) (int64, error) {
	k := r.genKey(key)
	return r.client.LLen(ctx, k).Result()
}

func (r *cache) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	k := r.genKey(key)
	return r.client.LRange(ctx, k, start, stop).Result()
}

func (r *cache) LPop(ctx context.Context, key string) (string, error) {
	k := r.genKey(key)
	result, err := r.client.LPop(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (r *cache) RPop(ctx context.Context, key string) (string, error) {
	k := r.genKey(key)
	result, err := r.client.RPop(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (r *cache) LIndex(ctx context.Context, key string, index int64) (string, error) {
	k := r.genKey(key)
	result, err := r.client.LIndex(ctx, k, index).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (r *cache) LSet(ctx context.Context, key string, index int64, value interface{}) error {
	k := r.genKey(key)
	return r.client.LSet(ctx, k, index, value).Err()
}

func (r *cache) LInsert(ctx context.Context, key, op string, pivot, value interface{}) (int64, error) {
	k := r.genKey(key)
	return r.client.LInsert(ctx, k, op, pivot, value).Result()
}

func (r *cache) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	k := r.genKey(key)
	return r.client.LRem(ctx, k, count, value).Result()
}

func (r *cache) LTrim(ctx context.Context, key string, start, stop int64) error {
	k := r.genKey(key)
	return r.client.LTrim(ctx, k, start, stop).Err()
}

func (r *cache) LPos(ctx context.Context, key string, value string) (int64, error) {
	k := r.genKey(key)
	result, err := r.client.LPos(ctx, k, value, goredis.LPosArgs{}).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (r *cache) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = r.genKey(keys[i])
	}
	result, err := r.client.BLPop(ctx, timeout, keys...).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (r *cache) BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = r.genKey(keys[i])
	}
	result, err := r.client.BRPop(ctx, timeout, keys...).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (r *cache) RPopLPush(ctx context.Context, source, destination string) (string, error) {
	source = r.genKey(source)
	destination = r.genKey(destination)
	result, err := r.client.RPopLPush(ctx, source, destination).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (r *cache) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error) {
	source = r.genKey(source)
	destination = r.genKey(destination)
	result, err := r.client.BRPopLPush(ctx, source, destination, timeout).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (r *cache) LMove(ctx context.Context, source, destination, srcpos, destpos string) (string, error) {
	source = r.genKey(source)
	destination = r.genKey(destination)
	result, err := r.client.LMove(ctx, source, destination, srcpos, destpos).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (r *cache) BLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) (string, error) {
	source = r.genKey(source)
	destination = r.genKey(destination)
	result, err := r.client.BLMove(ctx, source, destination, srcpos, destpos, timeout).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (r *cache) LPushX(ctx context.Context, key string, values ...interface{}) (int64, error) {
	k := r.genKey(key)
	return r.client.LPushX(ctx, k, values...).Result()
}

func (r *cache) RPushX(ctx context.Context, key string, values ...interface{}) (int64, error) {
	k := r.genKey(key)
	return r.client.RPushX(ctx, k, values...).Result()
}
