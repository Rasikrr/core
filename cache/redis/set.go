package redis

import (
	"context"
)

func (c *Client) SAdd(ctx context.Context, key string, members ...interface{}) error {
	key = c.genKey(key)
	return c.client.SAdd(ctx, key, members...).Err()
}

func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	key = c.genKey(key)
	return c.client.SMembers(ctx, key).Result()
}

func (c *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	key = c.genKey(key)
	return c.client.SIsMember(ctx, key, member).Result()
}

func (c *Client) SRem(ctx context.Context, key string, members ...interface{}) error {
	key = c.genKey(key)
	return c.client.SRem(ctx, key, members...).Err()
}

func (c *Client) SCard(ctx context.Context, key string) (int64, error) {
	key = c.genKey(key)
	return c.client.SCard(ctx, key).Result()
}

func (c *Client) SMove(ctx context.Context, source, destination string, member interface{}) (bool, error) {
	source = c.genKey(source)
	destination = c.genKey(destination)
	return c.client.SMove(ctx, source, destination, member).Result()
}

func (c *Client) SPop(ctx context.Context, key string) (string, error) {
	key = c.genKey(key)
	return c.client.SPop(ctx, key).Result()
}

func (c *Client) SPopN(ctx context.Context, key string, count int64) ([]string, error) {
	key = c.genKey(key)
	return c.client.SPopN(ctx, key, count).Result()
}

func (c *Client) SRandMember(ctx context.Context, key string) (string, error) {
	key = c.genKey(key)
	return c.client.SRandMember(ctx, key).Result()
}

func (c *Client) SRandMemberN(ctx context.Context, key string, count int64) ([]string, error) {
	key = c.genKey(key)
	return c.client.SRandMemberN(ctx, key, count).Result()
}

func (c *Client) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SDiff(ctx, keys...).Result()
}

func (c *Client) SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = c.genKey(destination)
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SDiffStore(ctx, destination, keys...).Result()
}

func (c *Client) SInter(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SInter(ctx, keys...).Result()
}

func (c *Client) SInterStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = c.genKey(destination)
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SInterStore(ctx, destination, keys...).Result()
}

func (c *Client) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SUnion(ctx, keys...).Result()
}

func (c *Client) SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = c.genKey(destination)
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	return c.client.SUnionStore(ctx, destination, keys...).Result()
}
