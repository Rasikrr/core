// nolint
package metrics

import (
	"errors"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricer *Metricer
	once     sync.Once
)

type Metricer struct {
	nameSpace  string
	enabled    bool
	registerer prometheus.Registerer
}

var ErrMetricAlreadyExists = errors.New("metric already registered")

func Init(active bool, namespace string, reg prometheus.Registerer) {
	once.Do(func() {
		r := prometheus.DefaultRegisterer
		if reg != nil {
			r = reg
		}
		metricer = &Metricer{
			nameSpace:  namespace,
			enabled:    active,
			registerer: r,
		}
	})
}

func Enabled() bool { return metricer != nil && metricer.enabled }

func NameSpace() string {
	if metricer == nil {
		return ""
	}
	return metricer.nameSpace
}

func registerOrExisting(c prometheus.Collector) (prometheus.Collector, error) {
	if metricer == nil || !metricer.enabled {
		// метрики выключены — не регистрируем
		return c, nil
	}
	if err := metricer.registerer.Register(c); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return are.ExistingCollector, ErrMetricAlreadyExists
		}
		return nil, err
	}
	return c, nil
}
