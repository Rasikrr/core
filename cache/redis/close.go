package redis

import (
	"context"

	"github.com/Rasikrr/core/log"
)

func (r *cache) Close(ctx context.Context) error {
	log.Info(ctx, "closing redis")
	return r.client.Close()
}
