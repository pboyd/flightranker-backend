package server

import (
	"sort"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/pboyd/flightranker-backend/backendC/store"
)

// flightStatsByAirlineRow is one row in a response from
// flightStatsByAirlineQuery.
type flightStatsByAirlineRow struct {
	Airline          string    `json:"airline"`
	Flights          int       `json:"totalFlights"`
	OnTimePercentage float64   `json:"onTimePercentage"`
	LastFlight       time.Time `json:"lastFlight"`
}

// flightStatsByAirlineQuery defines the flightStatsByAirline GraphQL query.
// The store instance is used when resolving the query.
func flightStatsByAirlineQuery(st *store.Store) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(
			graphql.NewObject(
				graphql.ObjectConfig{
					Name: "airlineFlightStats",
					Fields: graphql.Fields{
						"airline":          &graphql.Field{Type: graphql.String},
						"totalFlights":     &graphql.Field{Type: graphql.Int},
						"onTimePercentage": &graphql.Field{Type: graphql.Float},
						"lastFlight":       &graphql.Field{Type: graphql.DateTime},
					},
				},
			),
		),
		Args: graphql.FieldConfigArgument{
			"origin": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "airport IATA code (e.g. LAX)",
			},
			"destination": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "airport IATA code (e.g. LAX)",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			origin, _ := params.Args["origin"].(string)
			dest, _ := params.Args["destination"].(string)

			stats, err := st.FlightStats(
				params.Context,
				origin, dest,
				store.FlightStatsOpts{
					TimeGroup: store.GroupByAvailable,
				},
			)

			if err == store.ErrInvalidAirportCode {
				return nil, nil
			}

			outStats := make([]flightStatsByAirlineRow, 0, len(stats))
			for airline, set := range stats {
				outStats = append(outStats, flightStatsByAirlineRow{
					Airline:          airline,
					Flights:          set[0].Flights,
					LastFlight:       set[0].End,
					OnTimePercentage: set[0].OnTime(),
				})
			}

			// Rank airlines from best to worst by on-time percentage.
			sort.Slice(outStats, func(a, b int) bool {
				return outStats[a].OnTimePercentage > outStats[b].OnTimePercentage
			})

			return outStats, err
		},
	}
}
