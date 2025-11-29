package redis

import (
	"context"
	"fmt"
	coreCache "github.com/Rasikrr/core/cache"
	"net"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	client *redis.Client
	prefix string
}

func NewRedisCache(ctx context.Context, cfg Config, prefix string) (Cache, error) {
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

	return &cache{
		client: client,
		prefix: coreCache.PrefixKey(prefix),
	}, nil
}

func (r *cache) genKey(k string) string {
	return fmt.Sprintf("%s:%s", r.prefix, k)
}
