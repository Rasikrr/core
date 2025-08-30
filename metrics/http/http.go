package metrics

import (
	"github.com/Rasikrr/core/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

type Metrics struct {
	reqs     *prometheus.CounterVec
	duration *prometheus.HistogramVec
	errors   *prometheus.CounterVec
}

func NewHTTPMetrics(m metrics.Metricer) *Metrics {
	ss := m.Subsystem("http")
	return &Metrics{
		reqs:     ss.CounterVec("requests_total", "Total HTTP requests.", []string{"method", "path"}),
		duration: ss.HistogramVec("request_duration_seconds", "HTTP request duration.", []string{"method", "path"}, prometheus.DefBuckets),
		errors:   ss.CounterVec("errors_total", "Total HTTP handler errors.", []string{"method", "path", "code"}),
	}
}

// Middleware
func (h *Metrics) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := wrapWriter(w)
		next.ServeHTTP(ww, r)
		path := r.URL.Path
		method := r.Method

		h.reqs.WithLabelValues(method, path).Inc()
		h.duration.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
		if ww.status >= 400 {
			h.errors.WithLabelValues(method, path, http.StatusText(ww.status)).Inc()
		}
	})
}

type mwWriter struct {
	http.ResponseWriter
	status int
}

func wrapWriter(w http.ResponseWriter) *mwWriter {
	return &mwWriter{ResponseWriter: w, status: 200}
}
func (w *mwWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
