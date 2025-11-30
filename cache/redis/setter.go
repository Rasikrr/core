package redis

import (
	"context"
	"errors"
	"time"

	coreCache "github.com/Rasikrr/core/cache"
	goredis "github.com/redis/go-redis/v9"
)

// Set sets the value of a key
func (c *Cache) Set(ctx context.Context, key string, value any) error {
	k := c.genKey(key)
	return c.client.Set(ctx, k, value, 0).Err()
}

// SetWithExpiration sets the value of a key with expiration
func (c *Cache) SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error {
	k := c.genKey(key)
	return c.client.Set(ctx, k, value, expiration).Err()
}

// SetNX sets the value of a key only if it does not exist
// Returns true if the key was set, false if it already existed
func (c *Cache) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	k := c.genKey(key)
	return c.client.SetNX(ctx, k, value, expiration).Result()
}

// SetXX sets the value of a key only if it already exists
func (c *Cache) SetXX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	k := c.genKey(key)
	return c.client.SetXX(ctx, k, value, expiration).Result()
}

// GetSet atomically sets key to value and returns the old value
func (c *Cache) GetSet(ctx context.Context, key string, value any) (string, error) {
	k := c.genKey(key)
	result, err := c.client.GetSet(ctx, k, value).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

// GetDel atomically gets and deletes a key
func (c *Cache) GetDel(ctx context.Context, key string) (string, error) {
	k := c.genKey(key)
	result, err := c.client.GetDel(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

// GetEx gets the value and optionally sets expiration
func (c *Cache) GetEx(ctx context.Context, key string, expiration time.Duration) (string, error) {
	k := c.genKey(key)
	result, err := c.client.GetEx(ctx, k, expiration).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

// MSet sets multiple keys to multiple values
func (c *Cache) MSet(ctx context.Context, keyValues ...any) error {
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
func (c *Cache) Append(ctx context.Context, key, value string) (int64, error) {
	k := c.genKey(key)
	return c.client.Append(ctx, k, value).Result()
}
