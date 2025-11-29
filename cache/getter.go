package cache

import (
	"context"
	"time"
)

type Getter interface {
	Get(ctx context.Context, key string) (any, error)
	GetString(ctx context.Context, key string) (string, error)
	GetBytes(ctx context.Context, key string) ([]byte, error)
	GetBool(ctx context.Context, key string) (bool, error)
	GetInt(ctx context.Context, key string) (int, error)
	GetInt64(ctx context.Context, key string) (int64, error)
	GetFloat32(ctx context.Context, key string) (float32, error)
	GetFloat64(ctx context.Context, key string) (float64, error)
	GetTime(ctx context.Context, key string) (time.Time, error)
	MGet(ctx context.Context, keys ...string) ([]any, error)
}
