package metrics

import (
	"context"
	"errors"
	"time"

	"github.com/Rasikrr/core/log"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	promCounter struct{ c prometheus.Counter }
	promGauge   struct{ g prometheus.Gauge }
	promHist    struct{ o prometheus.Observer }
	promSumm    struct{ o prometheus.Observer }

	promCounterVec struct{ v *prometheus.CounterVec }
	promGaugeVec   struct{ v *prometheus.GaugeVec }
	promHistVec    struct{ v *prometheus.HistogramVec }
	promSummVec    struct{ v *prometheus.SummaryVec }
)

func (p promCounter) Inc()          { p.c.Inc() }
func (p promCounter) Add(v float64) { p.c.Add(v) }

func (p promGauge) Inc()              { p.g.Inc() }
func (p promGauge) Dec()              { p.g.Dec() }
func (p promGauge) Add(v float64)     { p.g.Add(v) }
func (p promGauge) Sub(v float64)     { p.g.Sub(v) }
func (p promGauge) Set(v float64)     { p.g.Set(v) }
func (p promGauge) SetToCurrentTime() { p.g.SetToCurrentTime() }

func (p promHist) Observe(v float64) { p.o.Observe(v) }

func (p promSumm) Observe(v float64) { p.o.Observe(v) }

func (p promCounterVec) WithLabelValues(lvs ...string) Counter {
	return promCounter{c: p.v.WithLabelValues(lvs...)}
}
func (p promGaugeVec) WithLabelValues(lvs ...string) Gauge {
	return promGauge{g: p.v.WithLabelValues(lvs...)}
}
func (p promHistVec) WithLabelValues(lvs ...string) Histogram {
	return promHist{o: p.v.WithLabelValues(lvs...)} // <- Observer
}
func (p promSummVec) WithLabelValues(lvs ...string) Summary {
	return promSumm{o: p.v.WithLabelValues(lvs...)} // <- Observer
}

/* ========== internals ========== */

func registerer() prometheus.Registerer {
	if metricer == nil || metricer.registerer == nil {
		return prometheus.DefaultRegisterer
	}
	return metricer.registerer
}
func namespace() string {
	if metricer == nil {
		return ""
	}
	return metricer.namespace
}
func isOn() bool { return metricer != nil && metricer.enabled }

// typed register-or-get helpers
func registerCounter(c prometheus.Counter) prometheus.Counter {
	if err := registerer().Register(c); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			if ex, ok := are.ExistingCollector.(prometheus.Counter); ok {
				return ex
			}
			log.Fatalf(context.Background(), "metrics error: %v", err)
		}
	}
	return c
}
func registerGauge(g prometheus.Gauge) prometheus.Gauge {
	if err := registerer().Register(g); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			if ex, ok := are.ExistingCollector.(prometheus.Gauge); ok {
				return ex
			}
			log.Fatalf(context.Background(), "metrics error: %v", err)
		}
	}
	return g
}
func registerHistogram(h prometheus.Histogram) prometheus.Histogram {
	if err := registerer().Register(h); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			if ex, ok := are.ExistingCollector.(prometheus.Histogram); ok {
				return ex
			}
			log.Fatalf(context.Background(), "metrics error: %v", err)
		}
	}
	return h
}
func registerSummary(s prometheus.Summary) prometheus.Summary {
	if err := registerer().Register(s); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			if ex, ok := are.ExistingCollector.(prometheus.Summary); ok {
				return ex
			}
			log.Fatalf(context.Background(), "metrics error: %v", err)
		}
	}
	return s
}
func registerCounterVec(v *prometheus.CounterVec) *prometheus.CounterVec {
	if err := registerer().Register(v); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			if ex, ok := are.ExistingCollector.(*prometheus.CounterVec); ok {
				return ex
			}
			log.Fatalf(context.Background(), "metrics error: %v", err)
		}
	}
	return v
}
func registerGaugeVec(v *prometheus.GaugeVec) *prometheus.GaugeVec {
	if err := registerer().Register(v); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			if ex, ok := are.ExistingCollector.(*prometheus.GaugeVec); ok {
				return ex
			}
			log.Fatalf(context.Background(), "metrics error: %v", err)
		}
	}
	return v
}
func registerHistVec(v *prometheus.HistogramVec) *prometheus.HistogramVec {
	if err := registerer().Register(v); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			if ex, ok := are.ExistingCollector.(*prometheus.HistogramVec); ok {
				return ex
			}
			log.Fatalf(context.Background(), "metrics error: %v", err)
		}
	}
	return v
}
func registerSummVec(v *prometheus.SummaryVec) *prometheus.SummaryVec {
	if err := registerer().Register(v); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			if ex, ok := are.ExistingCollector.(*prometheus.SummaryVec); ok {
				return ex
			}
			log.Fatalf(context.Background(), "metrics error: %v", err)
		}
	}
	return v
}

/* ========== constructors ========== */

func NewCounter(subsystem, name, help string, constLabels prometheus.Labels) Counter {
	if !isOn() {
		return noopCounter{}
	}
	c := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   namespace(),
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
	})
	c = registerCounter(c)
	return promCounter{c: c}
}

func NewCounterVec(subsystem, name, help string, labelNames []string, constLabels prometheus.Labels) CounterVec {
	if !isOn() {
		return noopCounterVec{}
	}
	v := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace(),
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
	}, labelNames)
	v = registerCounterVec(v)
	return promCounterVec{v: v}
}

func NewGauge(subsystem, name, help string, constLabels prometheus.Labels) Gauge {
	if !isOn() {
		return noopGauge{}
	}
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace(),
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
	})
	g = registerGauge(g)
	return promGauge{g: g}
}

func NewGaugeVec(subsystem, name, help string, labelNames []string, constLabels prometheus.Labels) GaugeVec {
	if !isOn() {
		return noopGaugeVec{}
	}
	v := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace(),
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
	}, labelNames)
	v = registerGaugeVec(v)
	return promGaugeVec{v: v}
}

func NewHistogram(subsystem, name, help string, buckets []float64, constLabels prometheus.Labels) Histogram {
	if !isOn() {
		return noopHistogram{}
	}
	if len(buckets) == 0 {
		buckets = prometheus.DefBuckets
	}
	h := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   namespace(),
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
		Buckets:     buckets,
	})
	h = registerHistogram(h)
	return promHist{o: h}
}

func NewHistogramVec(subsystem, name, help string, buckets []float64, labelNames []string, constLabels prometheus.Labels) HistogramVec {
	if !isOn() {
		return noopHistVec{}
	}
	if len(buckets) == 0 {
		buckets = prometheus.DefBuckets
	}
	v := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   namespace(),
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
		Buckets:     buckets,
	}, labelNames)
	v = registerHistVec(v)
	return promHistVec{v: v}
}

type SummaryOpts struct {
	Objectives map[float64]float64
	MaxAge     time.Duration
	AgeBuckets uint32
	BufCap     uint32
}

func NewSummary(subsystem, name, help string, opts SummaryOpts, constLabels prometheus.Labels) Summary {
	if !isOn() {
		return noopSummary{}
	}
	s := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:   namespace(),
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
		Objectives:  opts.Objectives,
		MaxAge:      opts.MaxAge,
		AgeBuckets:  opts.AgeBuckets,
		BufCap:      opts.BufCap,
	})
	s = registerSummary(s)
	return promSumm{o: s}
}

func NewSummaryVec(subsystem, name, help string, opts SummaryOpts, labelNames []string, constLabels prometheus.Labels) SummaryVec {
	if !isOn() {
		return noopSummaryVec{}
	}
	v := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   namespace(),
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
		Objectives:  opts.Objectives,
		MaxAge:      opts.MaxAge,
		AgeBuckets:  opts.AgeBuckets,
		BufCap:      opts.BufCap,
	}, labelNames)
	v = registerSummVec(v)
	return promSummVec{v: v}
}
