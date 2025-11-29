package redis

import "context"

func (r *cache) Exists(ctx context.Context, key string) (bool, error) {
	k := r.genKey(key)
	exists, err := r.client.Exists(ctx, k).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
