package mode

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/* ========== 指标 ========== */

var (
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "total http requests",
		},
		[]string{"handler", "method"},
	)

	HTTPErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "total http errors",
		},
		[]string{"handler"},
	)
)

/* ========== 初始化 ========== */

func InitPrometheus() {
	prometheus.MustRegister(HTTPRequests)
	prometheus.MustRegister(HTTPErrors)
}

/* ========== metrics出口 ========== */

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

/* ========== response捕获状态码 ========== */

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.HandlerFunc, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		HTTPRequests.WithLabelValues(name, r.Method).Inc()

		rec := &responseWriter{ResponseWriter: w, status: 200}

		next(rec, r)

		if rec.status >= 400 {
			HTTPErrors.WithLabelValues(name).Inc()
		}
	}
}
