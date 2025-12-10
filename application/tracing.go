package application

import (
	"context"
	"fmt"

	"github.com/Rasikrr/core/tracing"
)

func (a *App) initTracing(ctx context.Context) error {
	err := tracing.Init(ctx, a.Config().Tracing, a.Config().AppName, a.Config().Environment.String())
	if err != nil {
		return fmt.Errorf("failed to init tracing: %v", err)
	}
	return nil
}
