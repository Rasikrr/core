package environment

import (
	"sync"

	"github.com/Rasikrr/core/enum"
)

var (
	once sync.Once
	env  enum.Environment
)

func SetEnv(e enum.Environment) {
	once.Do(func() {
		env = e
	})
}

func GetEnv() enum.Environment {
	return env
}
