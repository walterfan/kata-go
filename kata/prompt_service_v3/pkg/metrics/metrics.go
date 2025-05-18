package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)
//go_goroutines is already registered by Prometheus by default, there's no need to re-register it in your code.
var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	RequestErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "Number of failed HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
)

func Register() {
	prometheus.MustRegister(RequestCount, RequestDuration, RequestErrors)
}