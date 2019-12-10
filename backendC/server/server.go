package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/pboyd/flightranker-backend/backendC/store"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Run starts an HTTP server on port 8080 using the handler from the Handler
// function.
//
// If it's unable to start the program exits with an error, otherwise the
// function never returns.
func Run() {
	http.Handle("/", Handler())
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler returns an http.Handler that responds to GraphQL queries for flight
// stats.
func Handler() http.Handler {
	corsAllowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")

	store := store.New()

	queries := graphql.Fields{
		"airport":              airportQuery(store),
		"airportList":          airportListQuery(store),
		"flightStatsByAirline": flightStatsByAirlineQuery(store),
		"dailyFlightStats":     dailyFlightStatsQuery(store),
		"monthlyFlightStats":   monthlyFlightStatsQuery(store),
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
	prometheus.Unregister(requests)
	prometheus.MustRegister(requests)

	errors := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "graphql",
		Subsystem: name,
		Name:      "errors",
	})
	prometheus.Unregister(errors)
	prometheus.MustRegister(errors)

	inflight := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "graphql",
		Subsystem: name,
		Name:      "inflight",
	})
	prometheus.Unregister(inflight)
	prometheus.MustRegister(inflight)

	responseTime := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "graphql",
		Subsystem: name,
		Name:      "response_time",
	})
	prometheus.Unregister(responseTime)
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
