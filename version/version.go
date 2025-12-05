package version

import (
	"sync"
)

var (
	versionOnce sync.Once
	version     string
)

func SetVersion(v string) {
	versionOnce.Do(func() {
		version = v
	})
}

func GetVersion() string {
	return version
}
