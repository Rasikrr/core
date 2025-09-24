package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxKey string

const (
	txCtxKey ctxKey = "postgres_tx"
)

type TXManager interface {
	Transaction(ctx context.Context, txOpts pgx.TxOptions, fn func(ctx context.Context) error) error
}

type txManager struct {
	pool *pgxpool.Pool
}

func NewTXManager(pool *pgxpool.Pool) TXManager {
	return &txManager{
		pool: pool,
	}
}

func (t *txManager) Transaction(ctx context.Context, txOpts pgx.TxOptions, fn func(ctx context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, txOpts)
	ctx = context.WithValue(ctx, txCtxKey, tx)

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	err = fn(ctx)
	return err
}
