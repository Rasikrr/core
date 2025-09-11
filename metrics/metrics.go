package metrics

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Rasikrr/core/util"
	"github.com/prometheus/client_golang/prometheus"
)

const metricsErrorFormat = "metric register conflict: %v (ns=%q, subsystem=%q, name=%q, labels=%v)"

var ErrMetricAlreadyExists = errors.New("metric already registered")

type Metricer struct {
	nameSpace  string
	enabled    bool
	registerer prometheus.Registerer
}

var (
	metricer *Metricer
	once     sync.Once
)

// Init инициализирует пакет. reg можно передать обёрнутый WrapRegistererWith(...), если нужно.
func Init(active bool, namespace string, reg prometheus.Registerer) {
	once.Do(func() {
		r := prometheus.DefaultRegisterer
		if !util.IsNil(reg) {
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

// ----------------- Внутренний helper регистрации -----------------

func registerOrExisting(c prometheus.Collector) (prometheus.Collector, error) {
	if metricer == nil {
		return nil, fmt.Errorf("metrics not initialized")
	}
	// Если метрики выключены — считаем noop по регистрации:
	// метрика создаётся, но не попадает в реестр -> не экспозится.
	if !metricer.enabled {
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

// ----------------- Counter -----------------

func NewCounter(name, subsystem, help string) (prometheus.Counter, error) {
	c := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: metricer.nameSpace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	})
	col, err := registerOrExisting(c)
	if err != nil && !errors.Is(err, ErrMetricAlreadyExists) {
		return nil, err
	}
	existing, ok := col.(prometheus.Counter)
	if !ok {
		return nil, fmt.Errorf(metricsErrorFormat, "type mismatch", metricer.nameSpace, subsystem, name, nil)
	}
	return existing, nil
}

func NewCounterVec(subsystem, name, help string, labels []string) (*prometheus.CounterVec, error) {
	cv := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricer.nameSpace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
		labels,
	)
	col, err := registerOrExisting(cv)
	if err != nil && !errors.Is(err, ErrMetricAlreadyExists) {
		return nil, err
	}
	existing, ok := col.(*prometheus.CounterVec)
	if !ok {
		return nil, fmt.Errorf(metricsErrorFormat, "type mismatch", metricer.nameSpace, subsystem, name, labels)
	}
	return existing, nil
}

// ----------------- Gauge -----------------

func NewGauge(name, subsystem, help string) (prometheus.Gauge, error) {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: metricer.nameSpace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	})
	col, err := registerOrExisting(g)
	if err != nil && !errors.Is(err, ErrMetricAlreadyExists) {
		return nil, err
	}
	existing, ok := col.(prometheus.Gauge)
	if !ok {
		return nil, fmt.Errorf(metricsErrorFormat, "type mismatch", metricer.nameSpace, subsystem, name, nil)
	}
	return existing, nil
}

func NewGaugeVec(subsystem, name, help string, labels []string) (*prometheus.GaugeVec, error) {
	gv := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricer.nameSpace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
		labels,
	)
	col, err := registerOrExisting(gv)
	if err != nil && !errors.Is(err, ErrMetricAlreadyExists) {
		return nil, err
	}
	existing, ok := col.(*prometheus.GaugeVec)
	if !ok {
		return nil, fmt.Errorf(metricsErrorFormat, "type mismatch", metricer.nameSpace, subsystem, name, labels)
	}
	return existing, nil
}

// ----------------- Histogram -----------------

func NewHistogram(name, subsystem, help string, buckets []float64) (prometheus.Histogram, error) {
	h := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: metricer.nameSpace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
		Buckets:   buckets, // можно передать prometheus.DefBuckets
	})
	col, err := registerOrExisting(h)
	if err != nil && !errors.Is(err, ErrMetricAlreadyExists) {
		return nil, err
	}
	existing, ok := col.(prometheus.Histogram)
	if !ok {
		return nil, fmt.Errorf(metricsErrorFormat, "type mismatch", metricer.nameSpace, subsystem, name, nil)
	}
	return existing, nil
}

func NewHistogramVec(subsystem, name, help string, buckets []float64, labels []string) (*prometheus.HistogramVec, error) {
	hv := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricer.nameSpace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		},
		labels,
	)
	col, err := registerOrExisting(hv)
	if err != nil && !errors.Is(err, ErrMetricAlreadyExists) {
		return nil, err
	}
	existing, ok := col.(*prometheus.HistogramVec)
	if !ok {
		return nil, fmt.Errorf(metricsErrorFormat, "type mismatch", metricer.nameSpace, subsystem, name, labels)
	}
	return existing, nil
}

// ----------------- Summary -----------------

func NewSummary(name, subsystem, help string, objectives map[float64]float64) (prometheus.Summary, error) {
	s := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:  metricer.nameSpace,
		Subsystem:  subsystem,
		Name:       name,
		Help:       help,
		Objectives: objectives, // можно nil -> без квантилей
	})
	col, err := registerOrExisting(s)
	if err != nil && !errors.Is(err, ErrMetricAlreadyExists) {
		return nil, err
	}
	existing, ok := col.(prometheus.Summary)
	if !ok {
		return nil, fmt.Errorf(metricsErrorFormat, "type mismatch", metricer.nameSpace, subsystem, name, nil)
	}
	return existing, nil
}

func NewSummaryVec(subsystem, name, help string, objectives map[float64]float64, labels []string) (*prometheus.SummaryVec, error) {
	sv := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  metricer.nameSpace,
			Subsystem:  subsystem,
			Name:       name,
			Help:       help,
			Objectives: objectives,
		},
		labels,
	)
	col, err := registerOrExisting(sv)
	if err != nil && !errors.Is(err, ErrMetricAlreadyExists) {
		return nil, err
	}
	existing, ok := col.(*prometheus.SummaryVec)
	if !ok {
		return nil, fmt.Errorf(metricsErrorFormat, "type mismatch", metricer.nameSpace, subsystem, name, labels)
	}
	return existing, nil
}
