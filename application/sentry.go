package application

import (
	"context"
	"time"

	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/sentry"
)

func (a *App) initSentry(ctx context.Context) error {
	if !a.Config().Sentry.Enabled {
		log.Info(ctx, "sentry disabled")
		return nil
	}
	err := sentry.Init(a.Config().Sentry, a.Config().AppName, a.Config().Env())
	if err != nil {
		return err
	}
	log.Info(ctx, "sentry initialized")
	return nil
}

func (a *App) flushSentry(ctx context.Context) {
	if !sentry.Enabled() {
		return
	}
	log.Info(ctx, "flushing sentry events before shutdown")
	if !sentry.Flush(5 * time.Second) {
		log.Warn(ctx, "not all sentry events were flushed")
	}
}
