package cache

import "context"

type Numeric interface {
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, delta int64) (int64, error)
	DecrBy(ctx context.Context, key string, delta int64) (int64, error)
	IncrByFloat(ctx context.Context, key string, delta float64) (float64, error)
	DecrByFloat(ctx context.Context, key string, delta float64) (float64, error)
}
