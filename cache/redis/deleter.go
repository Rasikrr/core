package redis

import "context"

func (c *Cache) Delete(ctx context.Context, key string) error {
	k := c.genKey(key)
	return c.client.Del(ctx, k).Err()
}

func (c *Cache) DeleteAll(ctx context.Context) error {
	return c.client.FlushDBAsync(ctx).Err()
}
