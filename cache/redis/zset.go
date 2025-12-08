// nolint: revive
package redis

import (
	"context"
)

// ZAdd adds members with scores to a sorted set
func (c *Client) ZAdd(ctx context.Context, key string, members ...Z) (int64, error) {
	key = c.genKey(key)
	return c.client.ZAdd(ctx, key, members...).Result()
}

// ZAddNX adds members only if they don't exist
func (c *Client) ZAddNX(ctx context.Context, key string, members ...Z) (int64, error) {
	key = c.genKey(key)
	return c.client.ZAddNX(ctx, key, members...).Result()
}

// ZAddXX updates members only if they exist
func (c *Client) ZAddXX(ctx context.Context, key string, members ...Z) (int64, error) {
	key = c.genKey(key)
	return c.client.ZAddXX(ctx, key, members...).Result()
}

// ZAddGT updates score only if new score is greater
func (c *Client) ZAddGT(ctx context.Context, key string, members ...Z) (int64, error) {
	key = c.genKey(key)
	return c.client.ZAddGT(ctx, key, members...).Result()
}

// ZAddLT updates score only if new score is less
func (c *Client) ZAddLT(ctx context.Context, key string, members ...Z) (int64, error) {
	key = c.genKey(key)
	return c.client.ZAddLT(ctx, key, members...).Result()
}

// ZRem removes members from a sorted set
func (c *Client) ZRem(ctx context.Context, key string, members ...any) (int64, error) {
	key = c.genKey(key)
	return c.client.ZRem(ctx, key, members...).Result()
}

// ZScore returns the score of a member
func (c *Client) ZScore(ctx context.Context, key, member string) (float64, error) {
	key = c.genKey(key)
	return c.client.ZScore(ctx, key, member).Result()
}

// ZCard returns the number of members in a sorted set
func (c *Client) ZCard(ctx context.Context, key string) (int64, error) {
	key = c.genKey(key)
	return c.client.ZCard(ctx, key).Result()
}

// ZCount returns the number of members with scores in the given range
func (c *Client) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	key = c.genKey(key)
	return c.client.ZCount(ctx, key, min, max).Result()
}

// ZIncrBy increments the score of a member
func (c *Client) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	key = c.genKey(key)
	return c.client.ZIncrBy(ctx, key, increment, member).Result()
}

// ZRange returns members by index range (ascending order)
// start and stop are 0-based indexes (0 = first, -1 = last)
func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	key = c.genKey(key)
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores returns members with scores by index range (ascending order)
func (c *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error) {
	key = c.genKey(key)
	result, err := c.client.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZRevRange returns members by index range (descending order)
func (c *Client) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	key = c.genKey(key)
	return c.client.ZRevRange(ctx, key, start, stop).Result()
}

// ZRevRangeWithScores returns members with scores by index range (descending order)
func (c *Client) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error) {
	key = c.genKey(key)
	result, err := c.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZRangeByScore returns members by score range (ascending order)
// min/max: use "-inf", "+inf" for infinity, or numeric values like "0", "100"
// Example: min="-inf", max="100" returns all members with score <= 100
func (c *Client) ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) ([]string, error) {
	key = c.genKey(key)
	return c.client.ZRangeByScore(ctx, key, opt).Result()
}

// ZRangeByScoreWithScores returns members with scores by score range
func (c *Client) ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) ([]Z, error) {
	key = c.genKey(key)
	result, err := c.client.ZRangeByScoreWithScores(ctx, key, opt).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZRevRangeByScore returns members by score range (descending order)
func (c *Client) ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) ([]string, error) {
	key = c.genKey(key)
	return c.client.ZRevRangeByScore(ctx, key, opt).Result()
}

// ZRevRangeByScoreWithScores returns members with scores by score range (descending order)
func (c *Client) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) ([]Z, error) {
	key = c.genKey(key)
	result, err := c.client.ZRevRangeByScoreWithScores(ctx, key, opt).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZRank returns the rank (index) of a member (ascending order, 0-based)
func (c *Client) ZRank(ctx context.Context, key, member string) (int64, error) {
	key = c.genKey(key)
	return c.client.ZRank(ctx, key, member).Result()
}

// ZRevRank returns the rank of a member (descending order, 0-based)
func (c *Client) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	key = c.genKey(key)
	return c.client.ZRevRank(ctx, key, member).Result()
}

// ZRemRangeByRank removes members by rank range
func (c *Client) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	key = c.genKey(key)
	return c.client.ZRemRangeByRank(ctx, key, start, stop).Result()
}

// ZRemRangeByScore removes members by score range
func (c *Client) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	key = c.genKey(key)
	return c.client.ZRemRangeByScore(ctx, key, min, max).Result()
}

// ZRemRangeByLex removes members by lexicographical range
func (c *Client) ZRemRangeByLex(ctx context.Context, key, min, max string) (int64, error) {
	key = c.genKey(key)
	return c.client.ZRemRangeByLex(ctx, key, min, max).Result()
}

// ZPopMin removes and returns members with the lowest scores
func (c *Client) ZPopMin(ctx context.Context, key string, count ...int64) ([]Z, error) {
	key = c.genKey(key)
	result, err := c.client.ZPopMin(ctx, key, count...).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZPopMax removes and returns members with the highest scores
func (c *Client) ZPopMax(ctx context.Context, key string, count ...int64) ([]Z, error) {
	key = c.genKey(key)
	result, err := c.client.ZPopMax(ctx, key, count...).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZUnionStore computes union of sorted sets and stores result
func (c *Client) ZUnionStore(ctx context.Context, dest string, store *ZStore) (int64, error) {
	dest = c.genKey(dest)
	// Add prefix to all keys in store
	prefixedKeys := make([]string, len(store.Keys))
	for i, k := range store.Keys {
		prefixedKeys[i] = c.genKey(k)
	}
	store.Keys = prefixedKeys
	return c.client.ZUnionStore(ctx, dest, store).Result()
}

// ZInterStore computes intersection of sorted sets and stores result
func (c *Client) ZInterStore(ctx context.Context, dest string, store *ZStore) (int64, error) {
	dest = c.genKey(dest)
	// Add prefix to all keys in store
	prefixedKeys := make([]string, len(store.Keys))
	for i, k := range store.Keys {
		prefixedKeys[i] = c.genKey(k)
	}
	store.Keys = prefixedKeys
	return c.client.ZInterStore(ctx, dest, store).Result()
}

// ZDiff returns the difference between the first sorted set and all successive sets
func (c *Client) ZDiff(ctx context.Context, keys ...string) ([]string, error) {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = c.genKey(k)
	}
	return c.client.ZDiff(ctx, prefixedKeys...).Result()
}

// ZDiffWithScores returns the difference with scores
func (c *Client) ZDiffWithScores(ctx context.Context, keys ...string) ([]Z, error) {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = c.genKey(k)
	}
	result, err := c.client.ZDiffWithScores(ctx, prefixedKeys...).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZScan iterates over members of a sorted set
func (c *Client) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	key = c.genKey(key)
	return c.client.ZScan(ctx, key, cursor, match, count).Result()
}
