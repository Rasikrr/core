package redis

import (
	coreCache "github.com/Rasikrr/core/cache"
	"github.com/Rasikrr/core/interfaces"
)

type Cache interface {
	interfaces.Closer

	coreCache.Getter
	coreCache.Setter
	coreCache.List
	coreCache.Set
	coreCache.Hash
	coreCache.Numeric
	coreCache.Exister
	coreCache.Deleter
	coreCache.Expirer
	// TODO: Pipeline, PubSub,
}
