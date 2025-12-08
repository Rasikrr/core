package redis

import (
	"context"
)

func (c *Client) Close(_ context.Context) error {
	return c.client.Close()
}
