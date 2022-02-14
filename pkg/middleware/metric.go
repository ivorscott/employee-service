package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/devpies/employee-service/pkg/web"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var totalRequestsMetric = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatusMetric = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDurationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

// Metric middleware implements prometheus metrics.
func Metric() web.Middleware {
	f := func(before web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			v, ok := r.Context().Value(web.KeyValues).(*web.Values)
			if !ok {
				return fmt.Errorf("web value missing from context")
			}

			err := before(w, r)

			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()

			// time request
			timer := prometheus.NewTimer(httpDurationMetric.WithLabelValues(path))

			statusCode := v.StatusCode

			// increase metric counters using Inc() (i.e., increment)
			responseStatusMetric.WithLabelValues(strconv.Itoa(statusCode)).Inc()
			totalRequestsMetric.WithLabelValues(path).Inc()

			timer.ObserveDuration()

			return err
		}
		return h
	}
	return f
}
