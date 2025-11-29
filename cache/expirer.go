package cache

import (
	"context"
	"time"
)

type Expirer interface {
	Expire(ctx context.Context, key string, expiration time.Duration) error
}
