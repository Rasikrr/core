package cache

import (
	"fmt"
	"strings"

	"github.com/Rasikrr/core/environment"
)

func PrefixKey(key string) string {
	return fmt.Sprintf("%s:%s", strings.ToUpper(key), strings.ToUpper(environment.GetEnv().String()))
}
