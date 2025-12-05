package redis

import (
	"context"
)

func (c *Cache) Close(_ context.Context) error {
	return c.client.Close()
}
