package main

import (
	"github.com/graphql-go/graphql"
	"github.com/prometheus/client_golang/prometheus"
)

func graphQLMetrics(name string, fn graphql.FieldResolveFn) graphql.FieldResolveFn {
	requests := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "graphql",
		Subsystem: name,
		Name:      "requests",
	})
	prometheus.MustRegister(requests)

	errors := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "graphql",
		Subsystem: name,
		Name:      "errors",
	})
	prometheus.MustRegister(errors)

	inflight := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "graphql",
		Subsystem: name,
		Name:      "inflight",
	})
	prometheus.MustRegister(inflight)

	responseTime := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "graphql",
		Subsystem: name,
		Name:      "response_time",
	})
	prometheus.MustRegister(responseTime)

	return func(p graphql.ResolveParams) (r interface{}, err error) {
		timer := prometheus.NewTimer(responseTime)
		requests.Inc()
		inflight.Inc()
		defer func() {
			timer.ObserveDuration()
			inflight.Dec()

			if err != nil {
				errors.Inc()
			}
		}()
		return fn(p)
	}
}
