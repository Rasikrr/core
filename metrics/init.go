// nolint
package metrics

import (
	"sync"

	"github.com/Rasikrr/core/util"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricer *Metricer
	once     sync.Once
)

type Metricer struct {
	namespace  string
	enabled    bool
	registerer prometheus.Registerer
}

func Init(active bool, namespace string, reg prometheus.Registerer) {
	once.Do(func() {
		r := prometheus.DefaultRegisterer
		if !util.IsNil(reg) {
			r = reg
		}
		metricer = &Metricer{
			namespace:  namespace,
			enabled:    active,
			registerer: r,
		}
	})
}

func Enabled() bool { return metricer != nil && metricer.enabled }
