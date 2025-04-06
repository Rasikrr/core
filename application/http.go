package application

import (
	"context"
	"github.com/Rasikrr/core/http"
	"log"
)

// nolint: unparam
func (a *App) initHTTP(ctx context.Context) error {
	if !a.Config().HTTPConfig().Required {
		return nil
	}
	a.httpServer = http.NewServer(ctx, a.Config().HTTPConfig())
	log.Println("http initialized")
	a.starters.Add(a.httpServer)
	a.closers.Add(a.httpServer)
	return nil
}
