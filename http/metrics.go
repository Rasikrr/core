package http

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Rasikrr/core/metrics"
)

var (
	once sync.Once
	m    *Metrics
)

type Metrics struct {
	reqTotal   metrics.CounterVec
	inflight   metrics.Gauge
	latencySec metrics.HistogramVec
	reqBytes   metrics.Histogram
	resBytes   metrics.Histogram
}

func initHTTPMetrics() {
	once.Do(func() {
		durBuckets := []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
		sizeBuckets := []float64{200, 500, 1 << 10, 2 << 10, 5 << 10, 10 << 10, 50 << 10, 100 << 10, 1 << 20, 5 << 20, 10 << 20}

		m = &Metrics{
			reqTotal:   metrics.NewCounterVec("http", "requests_total", "HTTP requests", []string{"method", "code"}, nil),
			inflight:   metrics.NewGauge("http", "inflight", "In-flight requests", nil),
			latencySec: metrics.NewHistogramVec("http", "request_seconds", "HTTP request latency", durBuckets, []string{"method", "code"}, nil),
			reqBytes:   metrics.NewHistogram("http", "request_bytes", "HTTP request size", sizeBuckets, nil),
			resBytes:   metrics.NewHistogram("http", "response_bytes", "HTTP response size", sizeBuckets, nil),
		}
	})
}

/* ================= std net/http middleware ================= */
type rw struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func (m *Metrics) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.inflight.Inc()
		defer m.inflight.Dec()

		wrap := &rw{ResponseWriter: w, status: 200}
		next.ServeHTTP(wrap, r)

		code := strconv.Itoa(wrap.status)
		m.reqTotal.WithLabelValues(r.Method, code).Inc()
		m.latencySec.WithLabelValues(r.Method, code).Observe(time.Since(start).Seconds())

		if r.ContentLength > 0 {
			m.reqBytes.Observe(float64(r.ContentLength))
		}
		m.resBytes.Observe(float64(wrap.bytes))
	})
}

func (w *rw) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *rw) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += int64(n)
	return n, err
}
