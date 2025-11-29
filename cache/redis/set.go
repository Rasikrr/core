package redis

import (
	"context"
	"errors"

	coreCache "github.com/Rasikrr/core/cache"
	goredis "github.com/redis/go-redis/v9"
)

func (r *cache) SAdd(ctx context.Context, key string, members ...interface{}) error {
	key = r.genKey(key)
	return r.client.SAdd(ctx, key, members...).Err()
}

func (r *cache) SMembers(ctx context.Context, key string) ([]string, error) {
	key = r.genKey(key)
	return r.client.SMembers(ctx, key).Result()
}

func (r *cache) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	key = r.genKey(key)
	return r.client.SIsMember(ctx, key, member).Result()
}

func (r *cache) SRem(ctx context.Context, key string, members ...interface{}) error {
	key = r.genKey(key)
	return r.client.SRem(ctx, key, members...).Err()
}

func (r *cache) SCard(ctx context.Context, key string) (int64, error) {
	key = r.genKey(key)
	return r.client.SCard(ctx, key).Result()
}

func (r *cache) SMove(ctx context.Context, source, destination string, member interface{}) (bool, error) {
	source = r.genKey(source)
	destination = r.genKey(destination)
	return r.client.SMove(ctx, source, destination, member).Result()
}

func (r *cache) SPop(ctx context.Context, key string) (string, error) {
	key = r.genKey(key)
	result, err := r.client.SPop(ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (r *cache) SPopN(ctx context.Context, key string, count int64) ([]string, error) {
	key = r.genKey(key)
	return r.client.SPopN(ctx, key, count).Result()
}

func (r *cache) SRandMember(ctx context.Context, key string) (string, error) {
	key = r.genKey(key)
	result, err := r.client.SRandMember(ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (r *cache) SRandMemberN(ctx context.Context, key string, count int64) ([]string, error) {
	key = r.genKey(key)
	return r.client.SRandMemberN(ctx, key, count).Result()
}

func (r *cache) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = r.genKey(keys[i])
	}
	return r.client.SDiff(ctx, keys...).Result()
}

func (r *cache) SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = r.genKey(destination)
	for i := range keys {
		keys[i] = r.genKey(keys[i])
	}
	return r.client.SDiffStore(ctx, destination, keys...).Result()
}

func (r *cache) SInter(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = r.genKey(keys[i])
	}
	return r.client.SInter(ctx, keys...).Result()
}

func (r *cache) SInterStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = r.genKey(destination)
	for i := range keys {
		keys[i] = r.genKey(keys[i])
	}
	return r.client.SInterStore(ctx, destination, keys...).Result()
}

func (r *cache) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = r.genKey(keys[i])
	}
	return r.client.SUnion(ctx, keys...).Result()
}

func (r *cache) SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = r.genKey(destination)
	for i := range keys {
		keys[i] = r.genKey(keys[i])
	}
	return r.client.SUnionStore(ctx, destination, keys...).Result()
}
