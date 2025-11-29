package cache

import (
	"context"
	"time"
)

type Setter interface {
	Set(ctx context.Context, key string, value any) error
	MSet(ctx context.Context, keyValues ...any) error
	SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error
}
