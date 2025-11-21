package application

import (
	"context"

	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/redis"
)

func (a *App) initRedis(ctx context.Context) error {
	if !a.config.Redis.Required {
		return nil
	}

	var err error
	a.redis, err = redis.NewRedisCache(ctx, a.Config().Redis, a.Config().Name())
	if err != nil {
		return err
	}

	log.Info(ctx, "redis initialized")

	a.closers.Add(a.redis)

	return nil
}
