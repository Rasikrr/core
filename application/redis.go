package application

import (
	"context"
	"github.com/Rasikrr/core/redis"
	"log"
)

func (a *App) initRedis(ctx context.Context) error {
	if !a.config.Redis.Required {
		return nil
	}
	var err error
	a.redis, err = redis.NewRedisCache(ctx, a.Config().RedisConfig(), a.Config().Name())
	if err != nil {
		return err
	}
	log.Println("redis initialized")
	a.closers.Add(a.redis)
	return nil
}
