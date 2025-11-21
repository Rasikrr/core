package application

import (
	"context"

	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/sentry"
)

func (a *App) initSentry(ctx context.Context) error {
	if !a.Config().Sentry.Enabled {
		log.Info(ctx, "sentry disabled")
		return nil
	}
	err := sentry.Init(a.Config().Sentry, a.Config().Env())
	if err != nil {
		return err
	}
	log.Info(ctx, "sentry initialized")
	return nil
}
