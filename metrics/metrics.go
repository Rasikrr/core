package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metricer interface {
	Registerer() prometheus.Registerer
	Namespace() string
	Subsystem(name string) Subsystem
}

type PromMetricer struct {
	ns  string
	reg prometheus.Registerer
}

func NewPrometheusMetricer(namespace string, reg prometheus.Registerer) *PromMetricer {
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}
	return &PromMetricer{ns: namespace, reg: reg}
}

func (m *PromMetricer) Registerer() prometheus.Registerer { return m.reg }
func (m *PromMetricer) Namespace() string                 { return m.ns }
func (m *PromMetricer) Subsystem(name string) Subsystem {
	return &promSubsystem{ns: m.Namespace(), sub: name, reg: m.Registerer()}
}
