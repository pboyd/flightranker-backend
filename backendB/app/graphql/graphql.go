package graphql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/pboyd/flightranker-backend/backendb/app"
	"github.com/prometheus/client_golang/prometheus"
)

type ProcessorConfig struct {
	AirportStore     app.AirportStore
	FlightStatsStore app.FlightStatsStore
}

type Processor struct {
	config ProcessorConfig
	schema graphql.Schema
}

func NewProcessor(config ProcessorConfig) *Processor {
	processor := &Processor{
		config: config,
	}

	processor.schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "Query",
					Fields: graphql.Fields{
						"airport":              processor.airportQuery(),
						"airportList":          processor.airportListQuery(),
						"flightStatsByAirline": processor.flightStatsByAirlineQuery(),
						"dailyFlightStats":     processor.dailyFlightStatsQuery(),
					},
				},
			),
		},
	)

	return processor
}

func (p *Processor) Do(ctx context.Context, query string) (string, error) {
	result := graphql.Do(graphql.Params{
		Context:       ctx,
		Schema:        p.schema,
		RequestString: query,
	})

	if result.HasErrors() {
		return "", QueryError{errors: result.Errors}
	}

	buf, err := json.Marshal(result.Data)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

type QueryError struct {
	errors []gqlerrors.FormattedError
}

func (qe QueryError) Error() string {
	buf, err := json.Marshal(qe.errors)
	if err != nil {
		return fmt.Sprintf("error formatting errors: %#v", qe.errors)
	}
	return string(buf)
}

func instrumentResolver(name string, fn graphql.FieldResolveFn) graphql.FieldResolveFn {
	requests := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "graphql",
		Subsystem: name,
		Name:      "requests",
	})
	promRegister(requests)

	errors := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "graphql",
		Subsystem: name,
		Name:      "errors",
	})
	promRegister(errors)

	inflight := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "graphql",
		Subsystem: name,
		Name:      "inflight",
	})
	promRegister(inflight)

	responseTime := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "graphql",
		Subsystem: name,
		Name:      "response_time",
	})
	promRegister(responseTime)

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

func promRegister(c prometheus.Collector) {
	prometheus.Unregister(c)
	prometheus.MustRegister(c)
}
