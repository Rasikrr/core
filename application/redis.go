package application

import (
	"context"
	"fmt"

	"github.com/Rasikrr/core/cache/redis"
	"github.com/Rasikrr/core/log"
)

func (a *App) initRedis(ctx context.Context) error {
	if !a.config.Redis.Required {
		return nil
	}

	var err error
	prefix := fmt.Sprintf("%s:%s", a.Config().Environment.String(), a.Config().Name())
	a.redis, err = redis.NewRedisCache(ctx, a.Config().Redis, prefix)
	if err != nil {
		return err
	}

	log.Info(ctx, "cache initialized")

	a.closers.Add(a.redis)

	return nil
}
