package application

import (
	"context"
	"github.com/Rasikrr/core/database"
	"log"
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
	log.Println("postgres initialized")
	a.closers.Add(a.postgres)
	return nil
}
