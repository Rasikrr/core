package application

import (
	"context"

	"github.com/Rasikrr/core/database/postgres"
	"github.com/Rasikrr/core/log"
)

func (a *App) initPostgres(ctx context.Context) error {
	if !a.Config().Postgres.Required {
		return nil
	}

	var err error
	a.postgres, err = postgres.NewPostgres(ctx, a.Config().Postgres)
	if err != nil {
		return err
	}
	a.postgresTXManager = postgres.NewTXManager(a.postgres.Pool())

	log.Info(ctx, "postgres initialized")

	a.closers.Add(a.postgres)

	return nil
}
