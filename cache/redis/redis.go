package redis

import (
	"context"
	"fmt"
	"net"

	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/tracing"
	"github.com/redis/go-redis/extra/redisotel/v9"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	logger log.Logger
	client *redis.Client
	prefix string
}

func NewRedisCache(ctx context.Context, cfg Config, prefix string) (*Client, error) {
	addr := net.JoinHostPort(
		cfg.Host,
		cfg.Port,
	)

	opt := &redis.Options{
		Addr:         addr,
		ClientName:   fmt.Sprintf("redis-%s", prefix),
		Username:     cfg.User,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdle,
		MaxIdleConns: cfg.MaxIdle,
		ReadTimeout:  cfg.ReadTimeout,
	}

	client := redis.NewClient(opt)

	if tracing.Enabled() {
		err := redisotel.InstrumentTracing(client)
		if err != nil {
			return nil, err
		}
	}

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{
		logger: log.With(
			log.String("system", "redis"),
		),
		client: client,
		prefix: prefix,
	}, nil
}

func (c *Client) genKey(k string) string {
	return fmt.Sprintf("%s:%s", c.prefix, k)
}
