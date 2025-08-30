package metrics

import "github.com/prometheus/client_golang/prometheus"

type Subsystem interface {
	CounterVec(name, help string, labels []string) *prometheus.CounterVec
	GaugeVec(name, help string, labels []string) *prometheus.GaugeVec
	HistogramVec(name, help string, labels []string, buckets []float64) *prometheus.HistogramVec
}

type promSubsystem struct {
	ns, sub string
	reg     prometheus.Registerer
}

func (s *promSubsystem) CounterVec(name, help string, labels []string) *prometheus.CounterVec {
	v := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: s.ns, Subsystem: s.sub, Name: name, Help: help,
	}, labels)
	s.reg.MustRegister(v)
	return v
}

func (s *promSubsystem) GaugeVec(name, help string, labels []string) *prometheus.GaugeVec {
	v := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: s.ns, Subsystem: s.sub, Name: name, Help: help,
	}, labels)
	s.reg.MustRegister(v)
	return v
}

func (s *promSubsystem) HistogramVec(name, help string, labels []string, buckets []float64) *prometheus.HistogramVec {
	v := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: s.ns, Subsystem: s.sub, Name: name, Help: help, Buckets: buckets,
	}, labels)
	s.reg.MustRegister(v)
	return v
}
