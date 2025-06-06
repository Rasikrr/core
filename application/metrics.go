package application

import (
	"context"
	"github.com/Rasikrr/core/http"
	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/metrics"
)

// nolint: unparam
func (a *App) initMetrics(ctx context.Context) error {
	if !a.config.Metrics.Enabled {
		return nil
	}
	cfg := a.config.MetricsConfig()

	a.metrics = metrics.NewPrometheusMetricer(ctx, cfg.Namespace)
	a.metricsServer = http.NewMetricsServer(ctx, cfg)

	a.starters.Add(a.metricsServer)
	a.closers.Add(a.metricsServer)

	log.Info(ctx, "metrics initialized")

	return nil
}
