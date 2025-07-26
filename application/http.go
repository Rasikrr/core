package application

import (
	"context"
	"github.com/Rasikrr/core/http"
	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/metrics"
)

// nolint: unparam
func (a *App) initHTTP(ctx context.Context) error {
	if !a.Config().HTTPConfig().Required {
		return nil
	}

	var httpMetrics metrics.HTTPMetrics
	if a.metrics != nil {
		httpMetrics = a.metrics.HTTP()
	}

	a.httpServer = http.NewServer(
		ctx,
		a.Config().HTTPConfig(),
		a.Config().MetricsConfig(),
		httpMetrics,
	)

	log.Info(ctx, "http initialized")

	a.starters.Add(a.httpServer)
	a.closers.Add(a.httpServer)

	return nil
}
