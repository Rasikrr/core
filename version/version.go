package version

import (
	"sync"

	"github.com/Rasikrr/core/enum"
)

var (
	versionOnce sync.Once
	version     enum.Environment
)

func SetVersion(v enum.Environment) {
	versionOnce.Do(func() {
		version = v
	})
}

func GetVersion() enum.Environment {
	return version
}
