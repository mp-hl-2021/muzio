package httpapi

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	totalHttpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_total_requests",
		Help: "Handled HTTP requests",
	}, []string{"code", "method"})
	inflightHttpRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_inflight_requests",
		Help: "HTTP requests inflight now",
	})
	durationHttpRequests = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_handlers_duration_seconds",
		Help: "HTTP requests handling duration",
	}, []string{"path"})
)

func measurer() func(next http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return promhttp.InstrumentHandlerCounter(
			totalHttpRequests,
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request){
					start := time.Now()
					inflightHttpRequests.Inc()
					next.ServeHTTP(w, r)
					inflightHttpRequests.Dec()
					route := mux.CurrentRoute(r)
					path, _ := route.GetPathTemplate()

					durationHttpRequests.WithLabelValues(path).Observe(time.Since(start).Seconds())
				}))
	}
}
