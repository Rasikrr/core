package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rasikrr/core/config"
	"github.com/Rasikrr/core/log"
	redis "github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"strings"
	"time"
)

type cache struct {
	client *redis.Client
	prefix string
}

func NewRedisCache(ctx context.Context, cfg config.RedisConfig, prefix string) (Cache, error) {
	addr := hostPort(
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
		prefix: prefixKey(prefix),
	}, nil
}

func (r *cache) Close(ctx context.Context) error {
	log.Info(ctx, "closing redis")
	return r.client.Close()
}

func (r *cache) Get(ctx context.Context, key string) (any, error) {
	k := r.getKey(key)
	val, err := r.redisStringCmd(ctx, k).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return val, nil
}

func (r *cache) MGet(ctx context.Context, keys ...string) ([]any, error) {
	keys = lo.Map(keys, func(k string, _ int) string {
		return r.getKey(k)
	})
	values, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (r *cache) Set(ctx context.Context, key string, value any) error {
	k := r.getKey(key)
	return r.client.Set(ctx, k, value, 0).Err()
}

func (r *cache) MSet(ctx context.Context, keyValues ...any) error {
	if len(keyValues)%2 != 0 {
		return errors.New("invalid keyValues: must be even")
	}
	for i := 0; i < len(keyValues); i += 2 {
		keyValues[i] = r.getKey(keyValues[i].(string))
	}
	return r.client.MSet(ctx, keyValues...).Err()
}

func (r *cache) SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error {
	k := r.getKey(key)
	return r.client.Set(ctx, k, value, expiration).Err()
}

func (r *cache) Exists(ctx context.Context, key string) (bool, error) {
	k := r.getKey(key)
	exists, err := r.client.Exists(ctx, k).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (r *cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	k := r.getKey(key)
	return r.client.Expire(ctx, k, expiration).Err()
}

func (r *cache) Delete(ctx context.Context, key string) error {
	k := r.getKey(key)
	return r.client.Del(ctx, k).Err()
}

func (r *cache) RPush(ctx context.Context, key string, value ...any) error {
	k := r.getKey(key)
	return r.client.RPush(ctx, k, value...).Err()
}

func (r *cache) LPush(ctx context.Context, key string, value ...any) error {
	k := r.getKey(key)
	return r.client.LPush(ctx, k, value...).Err()
}

func (r *cache) LLen(ctx context.Context, key string) (int64, error) {
	k := r.getKey(key)
	return r.client.LLen(ctx, k).Result()
}

func (r *cache) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	k := r.getKey(key)
	return r.client.LRange(ctx, k, start, stop).Result()
}

func (r *cache) Flush(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

func (r *cache) redisStringCmd(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

func (r *cache) getKey(k string) string {
	return fmt.Sprintf("%s:%s", r.prefix, k)
}

func hostPort(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func prefixKey(key string) string {
	return strings.ToUpper(key)
}
