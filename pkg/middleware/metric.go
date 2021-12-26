package middleware

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ivorscott/employee-service/pkg/web"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"net/http"
	"strconv"
)

func init() {
	// register metrics
	// required in order to expose metrics to the http handler
	if err := prometheus.Register(totalRequests); err != nil {
		log.Println("totalRequests metric failed")
	}
	if err := prometheus.Register(responseStatus); err != nil {
		log.Println("responseStatus metric failed")
	}
	if err := prometheus.Register(httpDuration); err != nil {
		log.Println("httpDuration metric failed")
	}
}

// metric 1
// create new counter metric
// counter metrics only increment values but can be reset to zero when the application restarts.
// the current value is not very relevant, values overtime are more significant.
var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

// metric 2
var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

// metric 3
var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

// Metrics middleware implements prometheus metrics.
func Metrics() web.Middleware {
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
			timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))

			statusCode := v.StatusCode

			// increase metric counters using Inc() (i.e., increment)
			responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
			totalRequests.WithLabelValues(path).Inc()

			timer.ObserveDuration()

			return err
		}
		return h
	}
	return f

}
