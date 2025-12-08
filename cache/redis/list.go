package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func (c *Client) RPush(ctx context.Context, key string, value ...any) error {
	k := c.genKey(key)
	return c.client.RPush(ctx, k, value...).Err()
}

func (c *Client) LPush(ctx context.Context, key string, value ...any) error {
	k := c.genKey(key)
	return c.client.LPush(ctx, k, value...).Err()
}

func (c *Client) LLen(ctx context.Context, key string) (int64, error) {
	k := c.genKey(key)
	return c.client.LLen(ctx, k).Result()
}

func (c *Client) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	k := c.genKey(key)
	return c.client.LRange(ctx, k, start, stop).Result()
}

func (c *Client) LPop(ctx context.Context, key string) (string, error) {
	k := c.genKey(key)
	result, err := c.client.LPop(ctx, k).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) RPop(ctx context.Context, key string) (string, error) {
	k := c.genKey(key)
	result, err := c.client.RPop(ctx, k).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) LIndex(ctx context.Context, key string, index int64) (string, error) {
	k := c.genKey(key)
	result, err := c.client.LIndex(ctx, k, index).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) LSet(ctx context.Context, key string, index int64, value interface{}) error {
	k := c.genKey(key)
	return c.client.LSet(ctx, k, index, value).Err()
}

func (c *Client) LInsert(ctx context.Context, key, op string, pivot, value interface{}) (int64, error) {
	k := c.genKey(key)
	return c.client.LInsert(ctx, k, op, pivot, value).Result()
}

func (c *Client) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	k := c.genKey(key)
	return c.client.LRem(ctx, k, count, value).Result()
}

func (c *Client) LTrim(ctx context.Context, key string, start, stop int64) error {
	k := c.genKey(key)
	return c.client.LTrim(ctx, k, start, stop).Err()
}

func (c *Client) LPos(ctx context.Context, key string, value string) (int64, error) {
	k := c.genKey(key)
	result, err := c.client.LPos(ctx, k, value, goredis.LPosArgs{}).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *Client) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	result, err := c.client.BLPop(ctx, timeout, keys...).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = c.genKey(keys[i])
	}
	result, err := c.client.BRPop(ctx, timeout, keys...).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) RPopLPush(ctx context.Context, source, destination string) (string, error) {
	source = c.genKey(source)
	destination = c.genKey(destination)
	result, err := c.client.RPopLPush(ctx, source, destination).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error) {
	source = c.genKey(source)
	destination = c.genKey(destination)
	result, err := c.client.BRPopLPush(ctx, source, destination, timeout).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) LMove(ctx context.Context, source, destination, srcpos, destpos string) (string, error) {
	source = c.genKey(source)
	destination = c.genKey(destination)
	result, err := c.client.LMove(ctx, source, destination, srcpos, destpos).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) BLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) (string, error) {
	source = c.genKey(source)
	destination = c.genKey(destination)
	result, err := c.client.BLMove(ctx, source, destination, srcpos, destpos, timeout).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) LPushX(ctx context.Context, key string, values ...interface{}) (int64, error) {
	k := c.genKey(key)
	return c.client.LPushX(ctx, k, values...).Result()
}

func (c *Client) RPushX(ctx context.Context, key string, values ...interface{}) (int64, error) {
	k := c.genKey(key)
	return c.client.RPushX(ctx, k, values...).Result()
}
