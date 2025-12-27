package database

import (
	"context"

	"github.com/Rasikrr/core/enum"
)

type TXManager interface {
	Transaction(ctx context.Context, txOpts TXOptions, fn func(ctx context.Context) error) error
}

type TXOptions struct {
	IsolationLevel enum.IsoLevel
	ReadOnly       bool
}
