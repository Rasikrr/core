package postgres

import (
	"context"

	"github.com/Rasikrr/core/database"
	"github.com/Rasikrr/core/enum"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxKey string

const (
	txCtxKey ctxKey = "core:postgres_tx"
)

var isoLevelMap = map[enum.IsoLevel]pgx.TxIsoLevel{
	enum.IsoLevelReadCommited:   pgx.ReadCommitted,
	enum.IsoLevelRepeatableRead: pgx.RepeatableRead,
	enum.IsoLevelSerializable:   pgx.Serializable,
}

func getPgxIsoLevel(level enum.IsoLevel) pgx.TxIsoLevel {
	if pgxIsoLevel, ok := isoLevelMap[level]; ok {
		return pgxIsoLevel
	}
	return pgx.ReadCommitted
}

type TXManager struct {
	pool *pgxpool.Pool
}

func NewTXManager(pool *pgxpool.Pool) *TXManager {
	return &TXManager{
		pool: pool,
	}
}

func (t *TXManager) Transaction(ctx context.Context, txOpts database.TXOptions, fn func(ctx context.Context) error) (err error) {
	if _, ok := ctx.Value(txCtxKey).(pgx.Tx); ok {
		return fn(ctx)
	}

	opts := pgx.TxOptions{
		IsoLevel:   getPgxIsoLevel(txOpts.IsolationLevel),
		AccessMode: pgx.ReadWrite,
	}
	if txOpts.ReadOnly {
		opts.AccessMode = pgx.ReadOnly
	}

	tx, err := t.pool.BeginTx(ctx, opts)
	if err != nil {
		return errors.WithStack(err)
	}
	txCtx := context.WithValue(ctx, txCtxKey, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.CombineErrors(err, errors.WithStack(rollbackErr))
			}
		} else {
			err = errors.WithStack(tx.Commit(ctx))
		}
	}()

	err = fn(txCtx)
	return err
}
