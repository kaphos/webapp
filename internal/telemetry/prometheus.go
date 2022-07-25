package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

var (
	// RequestsLatency tracks the total amount of time needed
	// to complete a request, from when the backend receives it
	// to when the user gets a response.
	RequestsLatency = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "kphs",
		Name:       "requests_latency_seconds",
		Help:       "The time taken from receiving a request to returning it",
		Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.025, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		MaxAge:     time.Hour * 24 * 21,
	}, []string{"method", "status"})

	// SQLLatency tracks the amount of time taken for each SQL query
	// to complete.
	SQLLatency = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "kphs",
		Name:       "sql_latency_seconds",
		Help:       "The time taken to perform a database query",
		Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.025, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		MaxAge:     time.Hour * 24 * 21,
	}, []string{"method"})

	// ErrCheckCount counts the number of times errchk
	// was used to check if there is an error. This can be used
	// in comparison with ErrCaughtCount to see the error rate
	// of the system.
	ErrCheckCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "kphs",
		Name:      "err_check_total",
		Help:      "The number of times that return values are checked if they contain an error",
	})

	// ErrCaughtCount counts the number of errors that have been caught
	// by the errorhandling package.
	ErrCaughtCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "kphs",
		Name:      "err_caught_total",
		Help:      "The number of times that return values checked were actually errors",
	})

	PromHandler = promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})
)
