package redis

import (
	"context"
	"fmt"
	"net"

	"github.com/Rasikrr/core/log"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	logger log.Logger
	client *redis.Client
	prefix string
}

func NewRedisCache(ctx context.Context, cfg Config, prefix string) (*Cache, error) {
	addr := net.JoinHostPort(
		cfg.Host,
		cfg.Port,
	)

	opt := &redis.Options{
		Addr:         addr,
		Username:     cfg.User,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdle,
		MaxIdleConns: cfg.MaxIdle,
		ReadTimeout:  cfg.ReadTimeout,
	}

	client := redis.NewClient(opt)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Cache{
		logger: log.With(log.String("system", "redis")),
		client: client,
		prefix: prefix,
	}, nil
}

func (c *Cache) genKey(k string) string {
	return fmt.Sprintf("%s:%s", c.prefix, k)
}
