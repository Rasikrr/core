package cache

import (
	"context"
	"time"
)

type Hash interface {
	HSet(ctx context.Context, key string, values ...interface{}) error
	HSetNX(ctx context.Context, key, field string, value interface{}) (bool, error)
	HGet(ctx context.Context, key, field string) (any, error)
	HGetBytes(ctx context.Context, key, field string) ([]byte, error)
	HGetString(ctx context.Context, key, field string) (string, error)
	HGetBool(ctx context.Context, key, field string) (bool, error)
	HGetInt(ctx context.Context, key, field string) (int, error)
	HGetInt64(ctx context.Context, key, field string) (int64, error)
	HGetFloat32(ctx context.Context, key, field string) (float32, error)
	HGetFloat64(ctx context.Context, key, field string) (float64, error)
	HGetTime(ctx context.Context, key, field string) (time.Time, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HMGet(ctx context.Context, key string, fields ...string) ([]any, error)
	HDel(ctx context.Context, key string, fields ...string) (int64, error)
	HExists(ctx context.Context, key, field string) (bool, error)
	HLen(ctx context.Context, key string) (int64, error)
	HKeys(ctx context.Context, key string) ([]string, error)
	HVals(ctx context.Context, key string) ([]string, error)
	HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error)
	HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error)
}
