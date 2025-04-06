package application

import (
	"github.com/Rasikrr/core/interfaces"
	"sync"
)

type Starters struct {
	mut      sync.Mutex
	starters []interfaces.Starter
}

func (s *Starters) Add(starter interfaces.Starter) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.starters = append(s.starters, starter)
}
