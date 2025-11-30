package redis

import (
	"context"
	"errors"
	"time"

	coreCache "github.com/Rasikrr/core/cache"
	goredis "github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

// Pipeliner wraps redis.Pipeliner with automatic key prefixing
type Pipeliner struct {
	pipe   goredis.Pipeliner
	prefix string
}

// Pipeline returns a new Pipeliner for manual control.
// Don't forget to call Exec() when done.
//
// Example:
//
//	pipe := cache.Pipeline()
//	pipe.Set(ctx, "key1", "value1", time.Hour)
//	pipe.Set(ctx, "key2", "value2", time.Hour)
//	if err := pipe.Exec(ctx); err != nil {
//	    return err
//	}
func (p *Pipeliner) Pipeline() *Pipeliner {
	return &Pipeliner{
		pipe:   p.pipe.Pipeline(),
		prefix: p.prefix,
	}
}

// PipelineExec executes multiple commands in a pipeline with automatic Exep.
// This is the recommended way for most use cases.
//
// Example:
//
//	err := cache.PipelineExec(ctx, func(pipe *Pipeliner) error {
//	    pipe.Set(ctx, "key1", "value1", time.Hour)
//	    pipe.Set(ctx, "key2", "value2", time.Hour)
//	    pipe.Incr(ctx, "counter")
//	    return nil
//	})
func (p *Pipeliner) PipelineExec(ctx context.Context, fn func(*Pipeliner) error) error {
	pipe := p.Pipeline()

	if err := fn(pipe); err != nil {
		return err
	}

	return pipe.Exec(ctx)
}

// Exec executes all queued commands
func (p *Pipeliner) Exec(ctx context.Context) error {
	_, err := p.pipe.Exec(ctx)
	return err
}

// Discard discards all queued commands
func (p *Pipeliner) Discard() {
	p.pipe.Discard()
}

// Len returns the number of queued commands
func (p *Pipeliner) Len() int {
	return p.pipe.Len()
}

func (p *Pipeliner) genKey(key string) string {
	return p.prefix + ":" + key
}

// ============================================================================
// Getter
// ============================================================================

func (p *Pipeliner) Get(ctx context.Context, key string) (any, error) {
	k := p.genKey(key)
	val, err := p.pipe.Get(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return val, nil
}

func (p *Pipeliner) GetString(ctx context.Context, key string) (string, error) {
	k := p.genKey(key)
	val, err := p.pipe.Get(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return val, nil
}

func (p *Pipeliner) GetBytes(ctx context.Context, key string) ([]byte, error) {
	k := p.genKey(key)
	val, err := p.pipe.Get(ctx, k).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return val, nil
}

func (p *Pipeliner) GetBool(ctx context.Context, key string) (bool, error) {
	k := p.genKey(key)
	val, err := p.pipe.Get(ctx, k).Bool()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return false, coreCache.ErrNotFound
		}
		return false, err
	}
	return val, nil
}

func (p *Pipeliner) GetInt(ctx context.Context, key string) (int, error) {
	k := p.genKey(key)
	val, err := p.pipe.Get(ctx, k).Int()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return val, nil
}

func (p *Pipeliner) GetInt64(ctx context.Context, key string) (int64, error) {
	k := p.genKey(key)
	val, err := p.pipe.Get(ctx, k).Int64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return val, nil
}

func (p *Pipeliner) GetFloat32(ctx context.Context, key string) (float32, error) {
	k := p.genKey(key)
	val, err := p.pipe.Get(ctx, k).Float32()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
	}
	return val, err
}
func (p *Pipeliner) GetFloat64(ctx context.Context, key string) (float64, error) {
	k := p.genKey(key)
	val, err := p.pipe.Get(ctx, k).Float64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
	}
	return val, err
}

func (p *Pipeliner) GetTime(ctx context.Context, key string) (time.Time, error) {
	k := p.genKey(key)
	val, err := p.pipe.Get(ctx, k).Time()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return time.Time{}, coreCache.ErrNotFound
		}
	}
	return val, err
}

func (p *Pipeliner) MGet(ctx context.Context, keys ...string) ([]any, error) {
	keys = lo.Map(keys, func(k string, _ int) string {
		return p.genKey(k)
	})
	values, err := p.pipe.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	return values, nil
}

// ============================================================================
// Deleter
// ============================================================================

func (p *Pipeliner) Delete(ctx context.Context, key string) error {
	k := p.genKey(key)
	return p.pipe.Del(ctx, k).Err()
}

func (p *Pipeliner) DeleteAll(ctx context.Context) error {
	return p.pipe.FlushDBAsync(ctx).Err()
}

// ============================================================================
// Exist
// ============================================================================

func (p *Pipeliner) Exists(ctx context.Context, key string) (bool, error) {
	k := p.genKey(key)
	exists, err := p.pipe.Exists(ctx, k).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

// ============================================================================
// Expire
// ============================================================================

// Expire sets a timeout on a key
func (p *Pipeliner) Expire(ctx context.Context, key string, expiration time.Duration) error {
	k := p.genKey(key)
	return p.pipe.Expire(ctx, k, expiration).Err()
}

// ExpireAt sets an absolute Unix timestamp expiration on a key
func (p *Pipeliner) ExpireAt(ctx context.Context, key string, tm time.Time) error {
	k := p.genKey(key)
	return p.pipe.ExpireAt(ctx, k, tm).Err()
}

// TTL returns the remaining time to live of a key
func (p *Pipeliner) TTL(ctx context.Context, key string) (time.Duration, error) {
	k := p.genKey(key)
	return p.pipe.TTL(ctx, k).Result()
}

// Persist removes the expiration from a key
func (p *Pipeliner) Persist(ctx context.Context, key string) error {
	k := p.genKey(key)
	return p.pipe.Persist(ctx, k).Err()
}

// ============================================================================
// Hash
// ============================================================================

func (p *Pipeliner) HSet(ctx context.Context, key string, values ...interface{}) error {
	key = p.genKey(key)
	return p.pipe.HSet(ctx, key, values...).Err()
}

func (p *Pipeliner) HSetNX(ctx context.Context, key, field string, value interface{}) (bool, error) {
	key = p.genKey(key)
	return p.pipe.HSetNX(ctx, key, field, value).Result()
}

func (p *Pipeliner) HGet(ctx context.Context, key, field string) (any, error) {
	key = p.genKey(key)
	result, err := p.pipe.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (p *Pipeliner) HGetBytes(ctx context.Context, key, field string) ([]byte, error) {
	key = p.genKey(key)
	result, err := p.pipe.HGet(ctx, key, field).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (p *Pipeliner) HGetString(ctx context.Context, key, field string) (string, error) {
	key = p.genKey(key)
	result, err := p.pipe.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (p *Pipeliner) HGetBool(ctx context.Context, key, field string) (bool, error) {
	key = p.genKey(key)
	result, err := p.pipe.HGet(ctx, key, field).Bool()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return false, coreCache.ErrNotFound
		}
		return false, err
	}
	return result, nil
}

func (p *Pipeliner) HGetInt(ctx context.Context, key, field string) (int, error) {
	key = p.genKey(key)
	result, err := p.pipe.HGet(ctx, key, field).Int()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (p *Pipeliner) HGetInt64(ctx context.Context, key, field string) (int64, error) {
	key = p.genKey(key)
	result, err := p.pipe.HGet(ctx, key, field).Int64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (p *Pipeliner) HGetFloat32(ctx context.Context, key, field string) (float32, error) {
	key = p.genKey(key)
	result, err := p.pipe.HGet(ctx, key, field).Float32()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (p *Pipeliner) HGetFloat64(ctx context.Context, key, field string) (float64, error) {
	key = p.genKey(key)
	result, err := p.pipe.HGet(ctx, key, field).Float64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (p *Pipeliner) HGetTime(ctx context.Context, key, field string) (time.Time, error) {
	key = p.genKey(key)
	result, err := p.pipe.HGet(ctx, key, field).Time()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return time.Time{}, coreCache.ErrNotFound
		}
		return time.Time{}, err
	}
	return result, nil
}

func (p *Pipeliner) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	key = p.genKey(key)
	return p.pipe.HGetAll(ctx, key).Result()
}

func (p *Pipeliner) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	key = p.genKey(key)
	return p.pipe.HMGet(ctx, key, fields...).Result()
}

func (p *Pipeliner) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	key = p.genKey(key)
	return p.pipe.HDel(ctx, key, fields...).Result()
}

func (p *Pipeliner) HExists(ctx context.Context, key, field string) (bool, error) {
	key = p.genKey(key)
	return p.pipe.HExists(ctx, key, field).Result()
}

func (p *Pipeliner) HLen(ctx context.Context, key string) (int64, error) {
	key = p.genKey(key)
	return p.pipe.HLen(ctx, key).Result()
}

func (p *Pipeliner) HKeys(ctx context.Context, key string) ([]string, error) {
	key = p.genKey(key)
	return p.pipe.HKeys(ctx, key).Result()
}

func (p *Pipeliner) HVals(ctx context.Context, key string) ([]string, error) {
	key = p.genKey(key)
	return p.pipe.HVals(ctx, key).Result()
}

func (p *Pipeliner) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	key = p.genKey(key)
	return p.pipe.HIncrBy(ctx, key, field, incr).Result()
}

func (p *Pipeliner) HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error) {
	key = p.genKey(key)
	return p.pipe.HIncrByFloat(ctx, key, field, incr).Result()
}

// ============================================================================
// List
// ============================================================================

func (p *Pipeliner) RPush(ctx context.Context, key string, value ...any) error {
	k := p.genKey(key)
	return p.pipe.RPush(ctx, k, value...).Err()
}

func (p *Pipeliner) LPush(ctx context.Context, key string, value ...any) error {
	k := p.genKey(key)
	return p.pipe.LPush(ctx, k, value...).Err()
}

func (p *Pipeliner) LLen(ctx context.Context, key string) (int64, error) {
	k := p.genKey(key)
	return p.pipe.LLen(ctx, k).Result()
}

func (p *Pipeliner) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	k := p.genKey(key)
	return p.pipe.LRange(ctx, k, start, stop).Result()
}

func (p *Pipeliner) LPop(ctx context.Context, key string) (string, error) {
	k := p.genKey(key)
	result, err := p.pipe.LPop(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (p *Pipeliner) RPop(ctx context.Context, key string) (string, error) {
	k := p.genKey(key)
	result, err := p.pipe.RPop(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (p *Pipeliner) LIndex(ctx context.Context, key string, index int64) (string, error) {
	k := p.genKey(key)
	result, err := p.pipe.LIndex(ctx, k, index).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (p *Pipeliner) LSet(ctx context.Context, key string, index int64, value interface{}) error {
	k := p.genKey(key)
	return p.pipe.LSet(ctx, k, index, value).Err()
}

func (p *Pipeliner) LInsert(ctx context.Context, key, op string, pivot, value interface{}) (int64, error) {
	k := p.genKey(key)
	return p.pipe.LInsert(ctx, k, op, pivot, value).Result()
}

func (p *Pipeliner) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	k := p.genKey(key)
	return p.pipe.LRem(ctx, k, count, value).Result()
}

func (p *Pipeliner) LTrim(ctx context.Context, key string, start, stop int64) error {
	k := p.genKey(key)
	return p.pipe.LTrim(ctx, k, start, stop).Err()
}

func (p *Pipeliner) LPos(ctx context.Context, key string, value string) (int64, error) {
	k := p.genKey(key)
	result, err := p.pipe.LPos(ctx, k, value, goredis.LPosArgs{}).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return result, nil
}

func (p *Pipeliner) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = p.genKey(keys[i])
	}
	result, err := p.pipe.BLPop(ctx, timeout, keys...).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (p *Pipeliner) BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = p.genKey(keys[i])
	}
	result, err := p.pipe.BRPop(ctx, timeout, keys...).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, coreCache.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}

func (p *Pipeliner) RPopLPush(ctx context.Context, source, destination string) (string, error) {
	source = p.genKey(source)
	destination = p.genKey(destination)
	result, err := p.pipe.RPopLPush(ctx, source, destination).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (p *Pipeliner) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error) {
	source = p.genKey(source)
	destination = p.genKey(destination)
	result, err := p.pipe.BRPopLPush(ctx, source, destination, timeout).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (p *Pipeliner) LMove(ctx context.Context, source, destination, srcpos, destpos string) (string, error) {
	source = p.genKey(source)
	destination = p.genKey(destination)
	result, err := p.pipe.LMove(ctx, source, destination, srcpos, destpos).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (p *Pipeliner) BLMove(ctx context.Context, source, destination, srcpos, destpos string, timeout time.Duration) (string, error) {
	source = p.genKey(source)
	destination = p.genKey(destination)
	result, err := p.pipe.BLMove(ctx, source, destination, srcpos, destpos, timeout).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (p *Pipeliner) LPushX(ctx context.Context, key string, values ...interface{}) (int64, error) {
	k := p.genKey(key)
	return p.pipe.LPushX(ctx, k, values...).Result()
}

func (p *Pipeliner) RPushX(ctx context.Context, key string, values ...interface{}) (int64, error) {
	k := p.genKey(key)
	return p.pipe.RPushX(ctx, k, values...).Result()
}

// ============================================================================
// Numeric
// ============================================================================

func (p *Pipeliner) Incr(ctx context.Context, key string) (int64, error) {
	key = p.genKey(key)
	return p.pipe.Incr(ctx, key).Result()
}

func (p *Pipeliner) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	key = p.genKey(key)
	return p.pipe.IncrBy(ctx, key, delta).Result()
}

func (p *Pipeliner) IncrByFloat(ctx context.Context, key string, delta float64) (float64, error) {
	key = p.genKey(key)
	return p.pipe.IncrByFloat(ctx, key, delta).Result()
}

func (p *Pipeliner) Decr(ctx context.Context, key string) (int64, error) {
	key = p.genKey(key)
	return p.pipe.Decr(ctx, key).Result()
}

func (p *Pipeliner) DecrBy(ctx context.Context, key string, delta int64) (int64, error) {
	key = p.genKey(key)
	return p.pipe.DecrBy(ctx, key, delta).Result()
}

func (p *Pipeliner) DecrByFloat(ctx context.Context, key string, delta float64) (float64, error) {
	key = p.genKey(key)
	return p.pipe.IncrByFloat(ctx, key, -delta).Result()
}

// ============================================================================
// Pub
// ============================================================================

func (p *Pipeliner) Publish(ctx context.Context, channel string, data any) error {
	return p.pipe.Publish(ctx, channel, data).Err()
}

// ============================================================================
// Scan
// ============================================================================

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
func (p *Pipeliner) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	prefixedMatch := p.genKey(match)

	keys, nextCursor, err := p.pipe.Scan(ctx, cursor, prefixedMatch, count).Result()
	if err != nil {
		return nil, 0, err
	}

	// Remove prefix from returned keys
	unprefixedKeys := make([]string, len(keys))
	prefixLen := len(p.prefix) + 1 // +1 for ":"
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
func (p *Pipeliner) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	key = p.genKey(key)
	return p.pipe.SScan(ctx, key, cursor, match, count).Result()
}

// HScan iterates over fields of a hash
func (p *Pipeliner) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	key = p.genKey(key)
	return p.pipe.HScan(ctx, key, cursor, match, count).Result()
}

// ============================================================================
// Set
// ============================================================================

func (p *Pipeliner) SAdd(ctx context.Context, key string, members ...interface{}) error {
	key = p.genKey(key)
	return p.pipe.SAdd(ctx, key, members...).Err()
}

func (p *Pipeliner) SMembers(ctx context.Context, key string) ([]string, error) {
	key = p.genKey(key)
	return p.pipe.SMembers(ctx, key).Result()
}

func (p *Pipeliner) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	key = p.genKey(key)
	return p.pipe.SIsMember(ctx, key, member).Result()
}

func (p *Pipeliner) SRem(ctx context.Context, key string, members ...interface{}) error {
	key = p.genKey(key)
	return p.pipe.SRem(ctx, key, members...).Err()
}

func (p *Pipeliner) SCard(ctx context.Context, key string) (int64, error) {
	key = p.genKey(key)
	return p.pipe.SCard(ctx, key).Result()
}

func (p *Pipeliner) SMove(ctx context.Context, source, destination string, member interface{}) (bool, error) {
	source = p.genKey(source)
	destination = p.genKey(destination)
	return p.pipe.SMove(ctx, source, destination, member).Result()
}

func (p *Pipeliner) SPop(ctx context.Context, key string) (string, error) {
	key = p.genKey(key)
	result, err := p.pipe.SPop(ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (p *Pipeliner) SPopN(ctx context.Context, key string, count int64) ([]string, error) {
	key = p.genKey(key)
	return p.pipe.SPopN(ctx, key, count).Result()
}

func (p *Pipeliner) SRandMember(ctx context.Context, key string) (string, error) {
	key = p.genKey(key)
	result, err := p.pipe.SRandMember(ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

func (p *Pipeliner) SRandMemberN(ctx context.Context, key string, count int64) ([]string, error) {
	key = p.genKey(key)
	return p.pipe.SRandMemberN(ctx, key, count).Result()
}

func (p *Pipeliner) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = p.genKey(keys[i])
	}
	return p.pipe.SDiff(ctx, keys...).Result()
}

func (p *Pipeliner) SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = p.genKey(destination)
	for i := range keys {
		keys[i] = p.genKey(keys[i])
	}
	return p.pipe.SDiffStore(ctx, destination, keys...).Result()
}

func (p *Pipeliner) SInter(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = p.genKey(keys[i])
	}
	return p.pipe.SInter(ctx, keys...).Result()
}

func (p *Pipeliner) SInterStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = p.genKey(destination)
	for i := range keys {
		keys[i] = p.genKey(keys[i])
	}
	return p.pipe.SInterStore(ctx, destination, keys...).Result()
}

func (p *Pipeliner) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	for i := range keys {
		keys[i] = p.genKey(keys[i])
	}
	return p.pipe.SUnion(ctx, keys...).Result()
}

func (p *Pipeliner) SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	destination = p.genKey(destination)
	for i := range keys {
		keys[i] = p.genKey(keys[i])
	}
	return p.pipe.SUnionStore(ctx, destination, keys...).Result()
}

// ============================================================================
// Setter
// ============================================================================

// Set sets the value of a key
func (p *Pipeliner) Set(ctx context.Context, key string, value any) error {
	k := p.genKey(key)
	return p.pipe.Set(ctx, k, value, 0).Err()
}

// SetWithExpiration sets the value of a key with expiration
func (p *Pipeliner) SetWithExpiration(ctx context.Context, key string, value any, expiration time.Duration) error {
	k := p.genKey(key)
	return p.pipe.Set(ctx, k, value, expiration).Err()
}

// SetNX sets the value of a key only if it does not exist
// Returns true if the key was set, false if it already existed
func (p *Pipeliner) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	k := p.genKey(key)
	return p.pipe.SetNX(ctx, k, value, expiration).Result()
}

// SetXX sets the value of a key only if it already exists
func (p *Pipeliner) SetXX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	k := p.genKey(key)
	return p.pipe.SetXX(ctx, k, value, expiration).Result()
}

// GetSet atomically sets key to value and returns the old value
func (p *Pipeliner) GetSet(ctx context.Context, key string, value any) (string, error) {
	k := p.genKey(key)
	result, err := p.pipe.GetSet(ctx, k, value).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

// GetDel atomically gets and deletes a key
func (p *Pipeliner) GetDel(ctx context.Context, key string) (string, error) {
	k := p.genKey(key)
	result, err := p.pipe.GetDel(ctx, k).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

// GetEx gets the value and optionally sets expiration
func (p *Pipeliner) GetEx(ctx context.Context, key string, expiration time.Duration) (string, error) {
	k := p.genKey(key)
	result, err := p.pipe.GetEx(ctx, k, expiration).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", coreCache.ErrNotFound
		}
		return "", err
	}
	return result, nil
}

// MSet sets multiple keys to multiple values
func (p *Pipeliner) MSet(ctx context.Context, keyValues ...any) error {
	if len(keyValues)%2 != 0 {
		return errors.New("invalid keyValues: must be even")
	}
	for i := 0; i < len(keyValues); i += 2 {
		keyStr, ok := keyValues[i].(string)
		if !ok {
			return errors.New("invalid key: must be string")
		}
		keyValues[i] = p.genKey(keyStr)
	}
	return p.pipe.MSet(ctx, keyValues...).Err()
}

// Append appends a value to a key
func (p *Pipeliner) Append(ctx context.Context, key, value string) (int64, error) {
	k := p.genKey(key)
	return p.pipe.Append(ctx, k, value).Result()
}

// ============================================================================
// Sorted Set (ZSet) operations
// ============================================================================

// ZAdd adds members with scores to a sorted set
func (p *Pipeliner) ZAdd(ctx context.Context, key string, members ...Z) (int64, error) {
	key = p.genKey(key)
	return p.pipe.ZAdd(ctx, key, members...).Result()
}

// ZAddNX adds members only if they don't exist
func (p *Pipeliner) ZAddNX(ctx context.Context, key string, members ...Z) (int64, error) {
	key = p.genKey(key)
	return p.pipe.ZAddNX(ctx, key, members...).Result()
}

// ZAddXX updates members only if they exist
func (p *Pipeliner) ZAddXX(ctx context.Context, key string, members ...Z) (int64, error) {
	key = p.genKey(key)
	return p.pipe.ZAddXX(ctx, key, members...).Result()
}

// ZAddGT updates score only if new score is greater
func (p *Pipeliner) ZAddGT(ctx context.Context, key string, members ...Z) (int64, error) {
	key = p.genKey(key)
	return p.pipe.ZAddGT(ctx, key, members...).Result()
}

// ZAddLT updates score only if new score is less
func (p *Pipeliner) ZAddLT(ctx context.Context, key string, members ...Z) (int64, error) {
	key = p.genKey(key)
	return p.pipe.ZAddLT(ctx, key, members...).Result()
}

// ZRem removes members from a sorted set
func (p *Pipeliner) ZRem(ctx context.Context, key string, members ...any) (int64, error) {
	key = p.genKey(key)
	return p.pipe.ZRem(ctx, key, members...).Result()
}

// ZScore returns the score of a member
func (p *Pipeliner) ZScore(ctx context.Context, key, member string) (float64, error) {
	key = p.genKey(key)
	score, err := p.pipe.ZScore(ctx, key, member).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return score, nil
}

// ZCard returns the number of members in a sorted set
func (p *Pipeliner) ZCard(ctx context.Context, key string) (int64, error) {
	key = p.genKey(key)
	return p.pipe.ZCard(ctx, key).Result()
}

// ZCount returns the number of members with scores in the given range
func (p *Pipeliner) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	key = p.genKey(key)
	return p.pipe.ZCount(ctx, key, min, max).Result()
}

// ZIncrBy increments the score of a member
func (p *Pipeliner) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	key = p.genKey(key)
	return p.pipe.ZIncrBy(ctx, key, increment, member).Result()
}

// ZRange returns members by index range (ascending order)
// start and stop are 0-based indexes (0 = first, -1 = last)
func (p *Pipeliner) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	key = p.genKey(key)
	return p.pipe.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores returns members with scores by index range (ascending order)
func (p *Pipeliner) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error) {
	key = p.genKey(key)
	result, err := p.pipe.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZRevRange returns members by index range (descending order)
func (p *Pipeliner) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	key = p.genKey(key)
	return p.pipe.ZRevRange(ctx, key, start, stop).Result()
}

// ZRevRangeWithScores returns members with scores by index range (descending order)
func (p *Pipeliner) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error) {
	key = p.genKey(key)
	result, err := p.pipe.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZRangeByScore returns members by score range (ascending order)
// min/max: use "-inf", "+inf" for infinity, or numeric values like "0", "100"
// Example: min="-inf", max="100" returns all members with score <= 100
func (p *Pipeliner) ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) ([]string, error) {
	key = p.genKey(key)
	return p.pipe.ZRangeByScore(ctx, key, opt).Result()
}

// ZRangeByScoreWithScores returns members with scores by score range
func (p *Pipeliner) ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) ([]Z, error) {
	key = p.genKey(key)
	result, err := p.pipe.ZRangeByScoreWithScores(ctx, key, opt).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZRevRangeByScore returns members by score range (descending order)
func (p *Pipeliner) ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) ([]string, error) {
	key = p.genKey(key)
	return p.pipe.ZRevRangeByScore(ctx, key, opt).Result()
}

// ZRevRangeByScoreWithScores returns members with scores by score range (descending order)
func (p *Pipeliner) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) ([]Z, error) {
	key = p.genKey(key)
	result, err := p.pipe.ZRevRangeByScoreWithScores(ctx, key, opt).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZRank returns the rank (index) of a member (ascending order, 0-based)
func (p *Pipeliner) ZRank(ctx context.Context, key, member string) (int64, error) {
	key = p.genKey(key)
	rank, err := p.pipe.ZRank(ctx, key, member).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return rank, nil
}

// ZRevRank returns the rank of a member (descending order, 0-based)
func (p *Pipeliner) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	key = p.genKey(key)
	rank, err := p.pipe.ZRevRank(ctx, key, member).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, coreCache.ErrNotFound
		}
		return 0, err
	}
	return rank, nil
}

// ZRemRangeByRank removes members by rank range
func (p *Pipeliner) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	key = p.genKey(key)
	return p.pipe.ZRemRangeByRank(ctx, key, start, stop).Result()
}

// ZRemRangeByScore removes members by score range
func (p *Pipeliner) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	key = p.genKey(key)
	return p.pipe.ZRemRangeByScore(ctx, key, min, max).Result()
}

// ZRemRangeByLex removes members by lexicographical range
func (p *Pipeliner) ZRemRangeByLex(ctx context.Context, key, min, max string) (int64, error) {
	key = p.genKey(key)
	return p.pipe.ZRemRangeByLex(ctx, key, min, max).Result()
}

// ZPopMin removes and returns members with the lowest scores
func (p *Pipeliner) ZPopMin(ctx context.Context, key string, count ...int64) ([]Z, error) {
	key = p.genKey(key)
	result, err := p.pipe.ZPopMin(ctx, key, count...).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZPopMax removes and returns members with the highest scores
func (p *Pipeliner) ZPopMax(ctx context.Context, key string, count ...int64) ([]Z, error) {
	key = p.genKey(key)
	result, err := p.pipe.ZPopMax(ctx, key, count...).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZUnionStore computes union of sorted sets and stores result
func (p *Pipeliner) ZUnionStore(ctx context.Context, dest string, store *ZStore) (int64, error) {
	dest = p.genKey(dest)
	// Add prefix to all keys in store
	prefixedKeys := make([]string, len(store.Keys))
	for i, k := range store.Keys {
		prefixedKeys[i] = p.genKey(k)
	}
	store.Keys = prefixedKeys
	return p.pipe.ZUnionStore(ctx, dest, store).Result()
}

// ZInterStore computes intersection of sorted sets and stores result
func (p *Pipeliner) ZInterStore(ctx context.Context, dest string, store *ZStore) (int64, error) {
	dest = p.genKey(dest)
	// Add prefix to all keys in store
	prefixedKeys := make([]string, len(store.Keys))
	for i, k := range store.Keys {
		prefixedKeys[i] = p.genKey(k)
	}
	store.Keys = prefixedKeys
	return p.pipe.ZInterStore(ctx, dest, store).Result()
}

// ZDiff returns the difference between the first sorted set and all successive sets
func (p *Pipeliner) ZDiff(ctx context.Context, keys ...string) ([]string, error) {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = p.genKey(k)
	}
	return p.pipe.ZDiff(ctx, prefixedKeys...).Result()
}

// ZDiffWithScores returns the difference with scores
func (p *Pipeliner) ZDiffWithScores(ctx context.Context, keys ...string) ([]Z, error) {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = p.genKey(k)
	}
	result, err := p.pipe.ZDiffWithScores(ctx, prefixedKeys...).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ZScan iterates over members of a sorted set
func (p *Pipeliner) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	key = p.genKey(key)
	return p.pipe.ZScan(ctx, key, cursor, match, count).Result()
}
