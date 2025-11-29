package redis

import "context"

func (r *cache) Incr(ctx context.Context, key string) (int64, error) {
	key = r.genKey(key)
	return r.client.Incr(ctx, key).Result()
}

func (r *cache) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	key = r.genKey(key)
	return r.client.IncrBy(ctx, key, delta).Result()
}

func (r *cache) IncrByFloat(ctx context.Context, key string, delta float64) (float64, error) {
	key = r.genKey(key)
	return r.client.IncrByFloat(ctx, key, delta).Result()
}

func (r *cache) Decr(ctx context.Context, key string) (int64, error) {
	key = r.genKey(key)
	return r.client.Decr(ctx, key).Result()
}

func (r *cache) DecrBy(ctx context.Context, key string, delta int64) (int64, error) {
	key = r.genKey(key)
	return r.client.DecrBy(ctx, key, delta).Result()
}

func (r *cache) DecrByFloat(ctx context.Context, key string, delta float64) (float64, error) {
	key = r.genKey(key)
	return r.client.IncrByFloat(ctx, key, -delta).Result()
}
