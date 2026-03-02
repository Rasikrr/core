package postgres

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
)

type txType struct{}

var (
	txCtxKey = txType{}
)

type TxManager struct {
	pool *Postgres
}

func NewTxManager(pool *Postgres) *TxManager {
	return &TxManager{pool: pool}
}

func (txm *TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	ok := HasTx(ctx)
	if ok {
		// Транзакция уже есть, просто выполняем функцию
		return fn(ctx)
	}

	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}
	tx, err := txm.pool.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	// ВАЖНО: defer выполняется в LIFO порядке, но recover должен быть первым (последним в коде),
	// чтобы поймать панику до того, как сработает логика commit/rollback, если вы их разделяете.
	// Но здесь у вас все в одной функции, это нормально.
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p) // Прокидываем панику дальше
		}

		// Логика завершения транзакции
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				// Комбинируем ошибку выполнения и ошибку роллбека
				err = errors.CombineErrors(err, rbErr)
			}
		} else {
			// Если ошибок не было, пробуем закоммитить
			// Поскольку 'err' именованный аргумент, это присвоение повлияет на return
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = errors.Wrap(commitErr, "failed to commit transaction")
			}
		}
	}()

	// Инжектим транзакцию в контекст
	txCtx := InjectTx(ctx, tx)

	err = fn(txCtx)
	return err
}

func HasTx(ctx context.Context) bool {
	_, ok := ctx.Value(txCtxKey).(pgx.Tx)
	return ok
}

func InjectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey, tx)
}
