package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelHandler = "handler"
	labelMethod  = "method"
	labelQuery   = "query"
	labelResult  = "result"
	labelService = "service"
	labelStatus  = "status"
)

var (
	Duration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "out_duration_seconds",
			Help:       "Summary of request duration in seconds",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{labelHandler, labelMethod, labelStatus},
	)
	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eout_rrors_total",
			Help: "Total number of errors",
		},
		[]string{labelHandler, labelMethod, labelStatus},
	)

	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "out_request_total",
			Help: "Total number of requests",
		},
		[]string{labelHandler, labelMethod},
	)

	DurationSql = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "out_duration_seconds",
			Help:       "Summary of request duration in seconds",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{labelQuery},
	)
	ErrorsTotalSql = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eout_rrors_total",
			Help: "Total number of errors",
		},
		[]string{labelQuery},
	)

	RequestsTotalSql = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "out_request_total",
			Help: "Total number of requests",
		},
		[]string{labelQuery},
	)
)

var (
	duration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "duration_seconds",
			Help:       "Summary of request duration in seconds",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{labelHandler, labelMethod, labelStatus},
	)

	errorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{labelHandler, labelMethod, labelStatus},
	)

	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_total",
			Help: "Total number of requests",
		},
		[]string{labelHandler, labelMethod},
	)
)

func init() {
	prometheus.MustRegister(duration)
	prometheus.MustRegister(errorsTotal)
	prometheus.MustRegister(requestsTotal)

	prometheus.MustRegister(Duration)
	prometheus.MustRegister(ErrorsTotal)
	prometheus.MustRegister(RequestsTotal)

	prometheus.MustRegister(DurationSql)
	prometheus.MustRegister(ErrorsTotalSql)
	prometheus.MustRegister(RequestsTotalSql)
}
