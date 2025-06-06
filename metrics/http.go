package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

type HTTPMetrics interface {
	IncRequest(method, path, status string)
	ObserveDuration(method, path, status string, duration float64)
}

type httpMetrics struct {
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func newHTTPMetrics(namespace string) HTTPMetrics {
	m := &httpMetrics{
		requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "Total HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "Duration of HTTP requests",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
	}

	prometheus.MustRegister(m.requests, m.duration)
	return m
}

func (m *httpMetrics) IncRequest(method, path, status string) {
	m.requests.WithLabelValues(method, path, status).Inc()
}

func (m *httpMetrics) ObserveDuration(method, path, status string, d float64) {
	m.duration.WithLabelValues(method, path, status).Observe(d)
}

type HTTPMetricsMiddleware struct {
	metrics HTTPMetrics
}

func NewHTTPMetricsMiddleware(metrics HTTPMetrics) *HTTPMetricsMiddleware {
	return &HTTPMetricsMiddleware{
		metrics: metrics,
	}
}

func (mid *HTTPMetricsMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds()
		method := r.Method
		path := r.URL.Path
		status := strconv.Itoa(rec.statusCode)

		mid.metrics.IncRequest(method, path, status)
		mid.metrics.ObserveDuration(method, path, status, duration)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
