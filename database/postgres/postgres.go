package postgres

import (
	"context"
	"time"

	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/tracing"
	"github.com/avast/retry-go"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

// nolint: gosec
func NewPostgres(ctx context.Context, cfg Config) (*Postgres, error) {
	conConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}
	conConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	conConfig.MaxConns = int32(cfg.MaxConns)
	conConfig.MinConns = int32(cfg.MinConns)
	conConfig.MaxConnIdleTime = cfg.MaxIdleConnIdleTime
	if tracing.Enabled() {
		conConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	}

	pool, err := pgxpool.NewWithConfig(ctx, conConfig)
	if err != nil {
		return nil, err
	}

	err = retry.Do(func() error {
		return pool.Ping(ctx)
	},
		retry.Attempts(3),
		retry.Delay(1*time.Second),
		retry.OnRetry(func(_ uint, err error) {
			log.Warnf(ctx, "Failed to connect to database: %v\n", err)
		}),
	)

	if err != nil {
		return nil, err
	}
	return &Postgres{
		pool: pool,
	}, nil
}

func (p *Postgres) Pool() *pgxpool.Pool {
	return p.pool
}

func (p *Postgres) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		log.Debugf(ctx, "Query: %s; elapsed: %v; args: %v\n", sql, elapsed, args)
	}()
	return p.GetQuerier(ctx).Query(ctx, sql, args...)
}

func (p *Postgres) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		log.Debugf(ctx, "Exec: %s; elapsed: %v; args: %v\n", sql, elapsed, args)
	}()
	return p.GetQuerier(ctx).Exec(ctx, sql, args...)
}

func (p *Postgres) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		log.Debugf(ctx, "QueryRow: %s; elapsed: %v; args: %v\n", sql, elapsed, args)
	}()
	return p.GetQuerier(ctx).QueryRow(ctx, sql, args...)
}

func (p *Postgres) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		log.Debugf(ctx, "BeginTx: %v; elapsed: %v\n", txOptions, elapsed)
	}()
	return p.pool.BeginTx(ctx, txOptions)
}

func (p *Postgres) Begin(ctx context.Context) (pgx.Tx, error) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		log.Debugf(ctx, "Begin: %v; elapsed: %v\n", nil, elapsed)
	}()
	return p.pool.Begin(ctx)
}

func (p *Postgres) Transaction(ctx context.Context, txOptions pgx.TxOptions, fn func(tx pgx.Tx) error) error {
	var err error
	tx, err := p.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()
	err = fn(tx)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		log.Debugf(ctx, "CopyFrom: %s; elapsed: %v; args: %v\n", tableName, elapsed, columnNames)
	}()
	return p.pool.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func (p *Postgres) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		log.Debugf(ctx, "SendBatch: %v; elapsed: %v\n", b, elapsed)
	}()
	return p.pool.SendBatch(ctx, b)
}

func (p *Postgres) Close(ctx context.Context) error {
	p.pool.Close()
	log.Info(ctx, "Postgres closed gracefully")
	return nil
}
