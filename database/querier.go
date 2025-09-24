package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func (p *Postgres) GetQuerier(ctx context.Context) Querier {
	if tx, ok := ctx.Value(txCtxKey).(pgx.Tx); ok {
		return tx
	}
	return p.pool
}
