package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics for monitoring with Prometheus.
var (
	GatherRuns = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gather_run_total",
			Help: "The number of times each gather command has been run",
		},
		[]string{"datatype"},
	)
	GatherErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gather_error_total",
			Help: "The number of times each gather command has had an error",
		},
		[]string{"datatype"},
	)
	GatherRuntime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "gather_command_runtime_seconds",
			Help: "How long each command took to run in seconds",
		},
		[]string{"datatype"},
	)
)

func init() {
	log.Println("Nodeinfo metrics have been initialized")
}
