package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/pboyd/flightranker-backend/backendC/store"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Run starts an HTTP server on port 8080 using the handler from the Handler
// function. If it's unable to start an error is returned, otherwise it never
// returns.
func Run(store *store.Store) error {
	http.Handle("/", Handler(store))
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":8080", nil)
}

// Handler returns an http.Handler that uses the store to respond to GraphQL
// queries.
func Handler(store *store.Store) http.Handler {
	corsAllowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")

	queries := graphql.Fields{
		"airport": newAirportQuery(store).Field(),
		//"airportList":          airportList,
		//"flightStatsByAirline": flightStatsByAirline,
		//"dailyFlightStats":     dailyFlightStats,
		//"monthlyFlightStats":   monthlyFlightStats,
	}

	// register each query with prometheus
	for key, query := range queries {
		instrumentResolver(key, query)
	}

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   "Query",
					Fields: queries,
				},
			),
		},
	)
	if err != nil {
		// This is a bug
		panic("server: failed to create graphql schema: " + err.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if corsAllowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
		}

		w.Header().Set("Content-Type", "application/json")

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("q"),
			Context:       r.Context(),
		})

		enc := json.NewEncoder(w)

		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(result.Errors)
			return
		}

		enc.Encode(result.Data)
	})
}

// instrumentResolver wraps the resolver function of a GraphQL query to record
// performance metrics in Prometheus.
//
// The Resolver function of the input graphql.Field will be modified.
func instrumentResolver(name string, query *graphql.Field) {
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

	fn := query.Resolve
	query.Resolve = func(p graphql.ResolveParams) (r interface{}, err error) {
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
