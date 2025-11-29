package application

import (
	"context"

	"github.com/Rasikrr/core/cache/redis"
	"github.com/Rasikrr/core/log"
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

	log.Info(ctx, "cache initialized")

	a.closers.Add(a.redis)

	return nil
}
