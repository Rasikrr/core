package interfaces

import (
	"context"
	"sync"
)

type Closers struct {
	mut     sync.Mutex
	closers []Closer
}

func NewClosers() *Closers {
	return &Closers{
		closers: make([]Closer, 0),
		mut:     sync.Mutex{},
	}
}

type Closer interface {
	Close(ctx context.Context) error
}

func (c *Closers) Add(closer Closer) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.closers = append(c.closers, closer)
}
