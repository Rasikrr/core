package application

import (
	"context"

	"github.com/Rasikrr/core/http"
	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/sentry"
)

// nolint: unparam
func (a *App) initSentry(ctx context.Context) error {
	if !a.Config().().Required {
		return nil
	}

	sentry.Init("TODO", a.Config().Environment)
}
