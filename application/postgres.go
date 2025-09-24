package application

import (
	"context"

	"github.com/Rasikrr/core/database"
	"github.com/Rasikrr/core/log"
)

func (a *App) initPostgres(ctx context.Context) error {
	if !a.config.Postgres.Required {
		return nil
	}

	var err error
	a.postgres, err = database.NewPostgres(ctx, a.Config().PostgresConfig())
	if err != nil {
		return err
	}
	a.postgresTXManager = database.NewTXManager(a.postgres.Pool())

	log.Info(ctx, "postgres initialized")

	a.closers.Add(a.postgres)

	return nil
}
