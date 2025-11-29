package redis

import "context"

func (r *cache) Delete(ctx context.Context, key string) error {
	k := r.genKey(key)
	return r.client.Del(ctx, k).Err()
}

func (r *cache) DeleteAll(ctx context.Context) error {
	return r.client.FlushDBAsync(ctx).Err()
}
