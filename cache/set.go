package cache

import (
	"context"
)

type Set interface {
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SRem(ctx context.Context, key string, members ...interface{}) error
	SCard(ctx context.Context, key string) (int64, error)
	SMove(ctx context.Context, source, destination string, member interface{}) (bool, error)
	SPop(ctx context.Context, key string) (string, error)
	SPopN(ctx context.Context, key string, count int64) ([]string, error)
	SRandMember(ctx context.Context, key string) (string, error)
	SRandMemberN(ctx context.Context, key string, count int64) ([]string, error)
	SDiff(ctx context.Context, keys ...string) ([]string, error)
	SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error)
	SInter(ctx context.Context, keys ...string) ([]string, error)
	SInterStore(ctx context.Context, destination string, keys ...string) (int64, error)
	SUnion(ctx context.Context, keys ...string) ([]string, error)
	SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error)
}
