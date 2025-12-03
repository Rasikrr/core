package redis

import (
	"github.com/redis/go-redis/v9"
)

type Z = redis.Z
type ZRangeBy = redis.ZRangeBy
type ZStore = redis.ZStore

type PubSub = redis.PubSub

type Message = redis.Message
type StringCMD = redis.StringCmd

var (
	Nil = redis.Nil
)
