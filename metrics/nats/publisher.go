package nats

import (
	"github.com/Rasikrr/core/metrics"
	"github.com/Rasikrr/core/util"
	"github.com/prometheus/client_golang/prometheus"
)

type PublisherMetrics interface {
	IncMsg(subject string)
	IncError(subject string)
	ObserveDuration(subject string, s float64)
}

type publisherMetrics struct {
	msgs    *prometheus.CounterVec
	errors  *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

func NewNATSPublisherMetrics(m metrics.Metricer) PublisherMetrics {
	if util.IsNil(m) {
		return &publisherNopMetrics{}
	}
	ss := m.Subsystem("nats_subscriber")
	return &publisherMetrics{
		msgs:    ss.CounterVec("messages_total", "Total NATS messages published.", []string{"subject"}),
		errors:  ss.CounterVec("errors_total", "Total NATS publish errors.", []string{"subject"}),
		latency: ss.HistogramVec("handle_duration_seconds", "NATS publish duration.", []string{"subject"}, prometheus.DefBuckets),
	}
}

func (n *publisherMetrics) IncMsg(subject string)   { n.msgs.WithLabelValues(subject).Inc() }
func (n *publisherMetrics) IncError(subject string) { n.errors.WithLabelValues(subject).Inc() }
func (n *publisherMetrics) ObserveDuration(subject string, s float64) {
	n.latency.WithLabelValues(subject).Observe(s)
}

type publisherNopMetrics struct{}

func (n *publisherNopMetrics) IncMsg(_ string)                     {}
func (n *publisherNopMetrics) IncError(_ string)                   {}
func (n *publisherNopMetrics) ObserveDuration(_ string, _ float64) {}
