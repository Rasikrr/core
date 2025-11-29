package cache

import (
	"context"
	"time"
)

type List interface {
	RPush(ctx context.Context, key string, value ...any) error
	LPush(ctx context.Context, key string, value ...any) error
	LLen(ctx context.Context, key string) (int64, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	LIndex(ctx context.Context, key string, index int64) (string, error)
	LSet(ctx context.Context, key string, index int64, value interface{}) error
	LInsert(ctx context.Context, key, op string, pivot, value interface{}) (int64, error)
	LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error)
	LTrim(ctx context.Context, key string, start, stop int64) error
	LPos(ctx context.Context, key string, value string) (int64, error)
	BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error)
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error)
	RPopLPush(ctx context.Context, source, destination string) (string, error)
	BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error)
	LMove(ctx context.Context, source, destination, srcpos, destpos string) (string, error)
	BLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) (string, error)
	LPushX(ctx context.Context, key string, values ...interface{}) (int64, error)
	RPushX(ctx context.Context, key string, values ...interface{}) (int64, error)
}
