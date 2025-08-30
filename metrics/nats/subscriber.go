package nats

import (
	"github.com/Rasikrr/core/metrics"
	"github.com/Rasikrr/core/util"
	"github.com/prometheus/client_golang/prometheus"
)

type SubscriberMetrics interface {
	IncMsg(subject string)
	IncError(subject string)
	ObserveDuration(subject string, s float64)
}

type subscriberMetrics struct {
	msgs    *prometheus.CounterVec
	errors  *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

func NewNATSSubscriberMetrics(m metrics.Metricer) SubscriberMetrics {
	if util.IsNil(m) {
		return &subscriberNopMetrics{}
	}
	ss := m.Subsystem("nats_subscriber")
	return &subscriberMetrics{
		msgs:    ss.CounterVec("messages_total", "Total NATS messages received.", []string{"subject"}),
		errors:  ss.CounterVec("errors_total", "Total NATS handler errors.", []string{"subject"}),
		latency: ss.HistogramVec("handle_duration_seconds", "NATS handler duration.", []string{"subject"}, prometheus.DefBuckets),
	}
}

func (n *subscriberMetrics) IncMsg(subject string)   { n.msgs.WithLabelValues(subject).Inc() }
func (n *subscriberMetrics) IncError(subject string) { n.errors.WithLabelValues(subject).Inc() }
func (n *subscriberMetrics) ObserveDuration(subject string, s float64) {
	n.latency.WithLabelValues(subject).Observe(s)
}

type subscriberNopMetrics struct{}

func (n *subscriberNopMetrics) IncMsg(_ string)                     {}
func (n *subscriberNopMetrics) IncError(_ string)                   {}
func (n *subscriberNopMetrics) ObserveDuration(_ string, _ float64) {}
