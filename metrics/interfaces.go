package metrics

type Counter interface {
	Inc()
	Add(float64)
}

type Gauge interface {
	Inc()
	Dec()
	Add(float64)
	Sub(float64)
	Set(float64)
	SetToCurrentTime()
}

type Histogram interface {
	Observe(float64)
}

type Summary interface {
	Observe(float64)
}

type CounterVec interface {
	WithLabelValues(...string) Counter
}
type GaugeVec interface {
	WithLabelValues(...string) Gauge
}
type HistogramVec interface {
	WithLabelValues(...string) Histogram
}
type SummaryVec interface {
	WithLabelValues(...string) Summary
}
