package redis

import (
	"context"
	"errors"
	"time"
)

// Set sets the value of a key
func (c *Client) Set(ctx context.Context, key string, value any) error {
	k := c.genKey(key)
	return c.client.Set(ctx, k, value, 0).Err()
}

// SetWithExpiration sets the value of a key with expiration
func (c *Client) SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error {
	k := c.genKey(key)
	return c.client.Set(ctx, k, value, expiration).Err()
}

// SetNX sets the value of a key only if it does not exist
// Returns true if the key was set, false if it already existed
func (c *Client) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	k := c.genKey(key)
	return c.client.SetNX(ctx, k, value, expiration).Result()
}

// SetXX sets the value of a key only if it already exists
func (c *Client) SetXX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	k := c.genKey(key)
	return c.client.SetXX(ctx, k, value, expiration).Result()
}

// MSet sets multiple keys to multiple values
func (c *Client) MSet(ctx context.Context, keyValues ...any) error {
	if len(keyValues)%2 != 0 {
		return errors.New("invalid keyValues: must be even")
	}
	for i := 0; i < len(keyValues); i += 2 {
		keyStr, ok := keyValues[i].(string)
		if !ok {
			return errors.New("invalid key: must be string")
		}
		keyValues[i] = c.genKey(keyStr)
	}
	return c.client.MSet(ctx, keyValues...).Err()
}

// Append appends a value to a key
func (c *Client) Append(ctx context.Context, key, value string) (int64, error) {
	k := c.genKey(key)
	return c.client.Append(ctx, k, value).Result()
}
