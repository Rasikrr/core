package application

import (
	"context"
	"fmt"

	"github.com/Rasikrr/core/brokers/nats"
	"github.com/Rasikrr/core/log"
)

func (a *App) initNats(ctx context.Context) error {
	cfg := a.Config().NATS
	if !cfg.Required {
		log.Debug(ctx, "nats is not required")
		return nil
	}

	var err error
	a.publisher, err = nats.NewPublisher(cfg.DSN)
	if err != nil {
		return fmt.Errorf("init NATS error: %w", err)
	}

	a.subscriber, err = nats.NewSubscriber(cfg.DSN, nats.WithQueue(cfg.Queue))
	if err != nil {
		return fmt.Errorf("init NATS error: %w", err)
	}

	log.Info(ctx, "nats initialized")

	a.starters.Add(a.subscriber)
	a.closers.Add(a.subscriber)

	a.closers.Add(a.publisher)

	return nil
}
