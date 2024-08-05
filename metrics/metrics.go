package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func init() {
	// Register the default Go runtime metrics
	prometheus.MustRegister(collectors.NewGoCollector())

	// Register process metrics (CPU, memory usage, etc.)
	prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}

var (
	TotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_requests_total",
		Help: "Total number of requests received by the API.",
	})

	ResponseDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "api_response_duration_seconds",
		Help:    "Histogram of response times for the API.",
		Buckets: prometheus.LinearBuckets(0.01, 0.05, 20),
	})
)
