package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	rps *prometheus.CounterVec
	ts  *prometheus.Histogram
}

func registerPrometheusMetrics() *metrics {
	hist := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      `ts`,
		Help:      `welcome handler response time`,
		Namespace: `welcome_handler`,
	})

	serverMetrics := &metrics{
		rps: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      `rps_total`,
			Help:      `RPS for welcome handler`,
			Namespace: `welcome_handler`,
		}, []string{`code`}),
		ts: &hist,
	}

	prometheus.MustRegister(serverMetrics.rps, *serverMetrics.ts)

	return serverMetrics
}
