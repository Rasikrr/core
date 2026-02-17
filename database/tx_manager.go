package database

import (
	"context"
)

type TXManager interface {
	Do(ctx context.Context, fn func(txCtx context.Context) error) error
}
