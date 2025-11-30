package redis

import "context"

// Scan iterates over keys matching a pattern
// WARNING: Use Scan instead of Keys in production to avoid blocking Redis
//
// Example:
//
//	cursor := uint64(0)
//	for {
//	    keys, nextCursor, err := cache.Scan(ctx, cursor, "user:*", 100)
//	    if err != nil {
//	        return err
//	    }
//	    // Process keys...
//	    cursor = nextCursor
//	    if cursor == 0 {
//	        break // iteration complete
//	    }
//	}
func (c *Cache) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	prefixedMatch := c.genKey(match)

	keys, nextCursor, err := c.client.Scan(ctx, cursor, prefixedMatch, count).Result()
	if err != nil {
		return nil, 0, err
	}

	// Remove prefix from returned keys
	unprefixedKeys := make([]string, len(keys))
	prefixLen := len(c.prefix) + 1 // +1 for ":"
	for i, key := range keys {
		if len(key) > prefixLen {
			unprefixedKeys[i] = key[prefixLen:]
		} else {
			unprefixedKeys[i] = key
		}
	}

	return unprefixedKeys, nextCursor, nil
}

// SScan iterates over members of a set
func (c *Cache) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	key = c.genKey(key)
	return c.client.SScan(ctx, key, cursor, match, count).Result()
}

// HScan iterates over fields of a hash
func (c *Cache) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	key = c.genKey(key)
	return c.client.HScan(ctx, key, cursor, match, count).Result()
}
