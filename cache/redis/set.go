package redis

import (
	"context"
	"errors"

	coreCache "github.com/Rasikrr/core/cache"
)

func (c *Cache) SAdd(ctx context.Context, key string, members ...interface{}) error {
	key = c.genKey(key)
	return c.client.SAdd(ctx, key, members...).Err()
}

func (c *Cache) SMembers(ctx context.Context, key string) ([]string, error) {
	key = c.genKey(key)
	return c.client.SMembers(ctx, key).Result()
}

func (c *Cache) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	key = c.genKey(key)
	return c.client.SIsMember(ctx, key, member).Result()
}

func (c *Cache) SRem(ctx context.Context, key string, members ...interface{}) error {
	key = c.genKey(key)
	return c.client.SRem(ctx, key, members...).Err()
}

func (c *Cache) SCard(ctx context.Context, key string) (int64, error) {
	key = c.genKey(key)
	return c.client.SCard(ctx, key).Result()
}

func (c *Cache) SMove(ctx context.Context, source, destination string, member interface{}) (bool, error) {
	source = c.genKey(source)
	destination = c.genKey(destination)
	return c.client.SMove(ctx, source, destination, member).Result()
}

func (c *Cache) SPop(ctx context.Context, key string) (string, error) {
	key = c.genKey(key)
	result, err := c.client.SPop(ctx, key).Result()
	if err != nil {
		if errors.Is(err, Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (c *Cache) SPopN(ctx context.Context, key string, count int64) ([]string, error) {
	key = c.genKey(key)
	return c.client.SPopN(ctx, key, count).Result()
}

func (c *Cache) SRandMember(ctx context.Context, key string) (string, error) {
	key = c.genKey(key)
	result, err := c.client.SRandMember(ctx, key).Result()
	if err != nil {
		if errors.Is(err, Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (c *Cache) SRandMemberN(ctx context.Context, key string, count int64) ([]string, error) {
	key = c.genKey(key)
	return c.client.SRandMemberN(ctx, key, count).Result()
}

func (c *Cache) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SDiff(ctx, keys...).Result()
}

func (c *Cache) SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = c.genKey(destination)
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SDiffStore(ctx, destination, keys...).Result()
}

func (c *Cache) SInter(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SInter(ctx, keys...).Result()
}

func (c *Cache) SInterStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = c.genKey(destination)
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SInterStore(ctx, destination, keys...).Result()
}

func (c *Cache) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SUnion(ctx, keys...).Result()
}

func (c *Cache) SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = c.genKey(destination)
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SUnionStore(ctx, destination, keys...).Result()
}
