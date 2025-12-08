package redis

import (
	"context"
	"time"
)

func (c *Client) HSet(ctx context.Context, key string, values ...interface{}) error {
	key = c.genKey(key)
	return c.client.HSet(ctx, key, values...).Err()
}

func (c *Client) HSetNX(ctx context.Context, key, field string, value interface{}) (bool, error) {
	key = c.genKey(key)
	return c.client.HSetNX(ctx, key, field, value).Result()
}

func (c *Client) HGet(ctx context.Context, key, field string) (any, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) HGetBytes(ctx context.Context, key, field string) ([]byte, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Bytes()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) HGetString(ctx context.Context, key, field string) (string, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) HGetBool(ctx context.Context, key, field string) (bool, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Bool()
	if err != nil {
		return false, err
	}
	return result, nil
}

func (c *Client) HGetInt(ctx context.Context, key, field string) (int, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Int()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *Client) HGetInt64(ctx context.Context, key, field string) (int64, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Int64()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *Client) HGetFloat32(ctx context.Context, key, field string) (float32, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Float32()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *Client) HGetFloat64(ctx context.Context, key, field string) (float64, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Float64()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (c *Client) HGetTime(ctx context.Context, key, field string) (time.Time, error) {
	key = c.genKey(key)
	result, err := c.client.HGet(ctx, key, field).Time()
	if err != nil {
		return time.Time{}, err
	}
	return result, nil
}

func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	key = c.genKey(key)
	return c.client.HGetAll(ctx, key).Result()
}

func (c *Client) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	key = c.genKey(key)
	return c.client.HMGet(ctx, key, fields...).Result()
}

func (c *Client) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	key = c.genKey(key)
	return c.client.HDel(ctx, key, fields...).Result()
}

func (c *Client) HExists(ctx context.Context, key, field string) (bool, error) {
	key = c.genKey(key)
	return c.client.HExists(ctx, key, field).Result()
}

func (c *Client) HLen(ctx context.Context, key string) (int64, error) {
	key = c.genKey(key)
	return c.client.HLen(ctx, key).Result()
}

func (c *Client) HKeys(ctx context.Context, key string) ([]string, error) {
	key = c.genKey(key)
	return c.client.HKeys(ctx, key).Result()
}

func (c *Client) HVals(ctx context.Context, key string) ([]string, error) {
	key = c.genKey(key)
	return c.client.HVals(ctx, key).Result()
}

func (c *Client) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	key = c.genKey(key)
	return c.client.HIncrBy(ctx, key, field, incr).Result()
}

func (c *Client) HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error) {
	key = c.genKey(key)
	return c.client.HIncrByFloat(ctx, key, field, incr).Result()
}
