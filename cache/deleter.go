package cache

import "context"

type Deleter interface {
	Delete(ctx context.Context, key string) error
	DeleteAll(ctx context.Context) error
}
