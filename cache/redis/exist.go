package redis

import "context"

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	k := c.genKey(key)
	exists, err := c.client.Exists(ctx, k).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
