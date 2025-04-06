package application

import (
	"github.com/Rasikrr/core/interfaces"
	"sync"
)

type Closers struct {
	mut     sync.Mutex
	closers []interfaces.Closer
}

func (c *Closers) Add(closer interfaces.Closer) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.closers = append(c.closers, closer)
}
