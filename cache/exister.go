package cache

import "context"

type Exister interface {
	Exists(ctx context.Context, key string) (bool, error)
}
