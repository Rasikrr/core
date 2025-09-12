package metrics

type (
	noopCounter    struct{}
	noopGauge      struct{}
	noopHistogram  struct{}
	noopSummary    struct{}
	noopCounterVec struct{}
	noopGaugeVec   struct{}
	noopHistVec    struct{}
	noopSummaryVec struct{}
)

func (noopCounter) Inc()                                 {}
func (noopCounter) Add(float64)                          {}
func (noopGauge) Inc()                                   {}
func (noopGauge) Dec()                                   {}
func (noopGauge) Add(float64)                            {}
func (noopGauge) Sub(float64)                            {}
func (noopGauge) Set(float64)                            {}
func (noopGauge) SetToCurrentTime()                      {}
func (noopHistogram) Observe(float64)                    {}
func (noopSummary) Observe(float64)                      {}
func (noopCounterVec) WithLabelValues(...string) Counter { return noopCounter{} }
func (noopGaugeVec) WithLabelValues(...string) Gauge     { return noopGauge{} }
func (noopHistVec) WithLabelValues(...string) Histogram  { return noopHistogram{} }
func (noopSummaryVec) WithLabelValues(...string) Summary { return noopSummary{} }
