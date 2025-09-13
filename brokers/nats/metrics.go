package nats

import (
	"sync"

	coreMetrics "github.com/Rasikrr/core/metrics"
)

type Metrics struct {
	// publish
	pubTotal coreMetrics.CounterVec   // {subject}
	pubBytes coreMetrics.HistogramVec // {subject}
	// receive (async handlers)
	recvTotal      coreMetrics.CounterVec   // {subject}
	recvBytes      coreMetrics.HistogramVec // {subject}
	handlerSeconds coreMetrics.HistogramVec // {subject}

	// request/reply
	reqTotal    coreMetrics.CounterVec   // {subject, outcome=ok|timeout|error}
	reqLatency  coreMetrics.HistogramVec // {subject, outcome}
	inflightReq coreMetrics.GaugeVec     // {subject}
}

var (
	metrics *Metrics
	once    sync.Once
)

func initNATSMetrics() {
	once.Do(func() {
		dur := []float64{0.0005, 0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
		sz := []float64{128, 256, 512, 1 << 10, 2 << 10, 4 << 10, 8 << 10, 16 << 10, 32 << 10, 64 << 10, 128 << 10, 1 << 20, 5 << 20, 10 << 20}

		metrics = &Metrics{
			pubTotal:       coreMetrics.NewCounterVec("nats", "publish_total", "NATS messages published", []string{"subject"}, nil),
			pubBytes:       coreMetrics.NewHistogramVec("nats", "publish_bytes", "Published message size (bytes)", sz, []string{"subject"}, nil),
			recvTotal:      coreMetrics.NewCounterVec("nats", "receive_total", "NATS messages received", []string{"subject"}, nil),
			recvBytes:      coreMetrics.NewHistogramVec("nats", "receive_bytes", "Received message size (bytes)", sz, []string{"subject"}, nil),
			handlerSeconds: coreMetrics.NewHistogramVec("nats", "handler_seconds", "Async handler duration (seconds)", dur, []string{"subject"}, nil),
			reqTotal:       coreMetrics.NewCounterVec("nats", "request_total", "NATS request calls", []string{"subject", "outcome"}, nil),
			reqLatency:     coreMetrics.NewHistogramVec("nats", "request_seconds", "NATS request latency (seconds)", dur, []string{"subject", "outcome"}, nil),
			inflightReq:    coreMetrics.NewGaugeVec("nats", "inflight_requests", "In-flight NATS requests", []string{"subject"}, nil),
		}
	})
}
