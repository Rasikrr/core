package application

import (
	"context"

	"github.com/Rasikrr/core/http"
	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/metrics"
)

// nolint: unparam
func (a *App) initMetrics(ctx context.Context) error {
	metrics.Init(
		a.Config().MetricsConfig().Enabled,
		a.Config().MetricsConfig().Namespace,
		nil,
	)
	if metrics.Enabled() {
		a.metricsServer = http.NewMetricsServer(ctx, a.Config().MetricsConfig())
		a.starters.Add(a.metricsServer)
		a.closers.Add(a.metricsServer)
		log.Infof(ctx, "metrics server initialized")
		log.Info(ctx, "metric initialized")
	}
	return nil
}
