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
			graphql.NewObject(graphql.ObjectConfig{
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
			"origin":      airportCodeArgument,
			"destination": airportCodeArgument,
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
			} else if err != nil {
				return nil, err
			}

			outStats := make([]flightStatsByAirlineRow, 0, len(stats))
			for _, airlineStats := range stats {
				outStats = append(outStats, flightStatsByAirlineRow{
					Airline:          airlineStats.Airline,
					Flights:          airlineStats.Rows[0].Flights,
					LastFlight:       airlineStats.Rows[0].End,
					OnTimePercentage: airlineStats.Rows[0].OnTime(),
				})
			}

			// Rank airlines from best to worst by on-time percentage.
			sort.Slice(outStats, func(a, b int) bool {
				return outStats[a].OnTimePercentage > outStats[b].OnTimePercentage
			})

			return outStats, nil
		},
	}
}

// flightStatsByDateType is the GraphQL definition of the return value from
// dailyFlightStatsQuery and monthylyFlightStatsQuery.
var flightStatsByDateType = graphql.NewList(
	graphql.NewObject(graphql.ObjectConfig{
		Name: "flightStatsByDate",
		Fields: graphql.Fields{
			"airline": &graphql.Field{Type: graphql.String},
			"rows": &graphql.Field{Type: graphql.NewList(graphql.NewObject(
				graphql.ObjectConfig{
					Name: "flightStatsByDateRow",
					Fields: graphql.Fields{
						"date":    &graphql.Field{Type: graphql.DateTime},
						"flights": &graphql.Field{Type: graphql.Int},
						"delays":  &graphql.Field{Type: graphql.Int},
						"onTimePercentage": &graphql.Field{
							Type:    graphql.Float,
							Resolve: resolveOnTimePercentage,
						},
					},
				},
			),
			)},
		},
	},
	),
)

// resolveOnTimePercentage is a graphql.Resolver that returns the result of the
// OnTime function from a source.StatsRow.
func resolveOnTimePercentage(params graphql.ResolveParams) (interface{}, error) {
	row, ok := params.Source.(store.StatsRow)
	if !ok {
		return 0, nil
	}

	return row.OnTime(), nil
}

// dailyFlightStatsQuery defines the dailyFlightStats GraphQL query, which
// returns flight stats grouped by day.
// The store instance is used when resolving the query.
func dailyFlightStatsQuery(st *store.Store) *graphql.Field {
	return &graphql.Field{
		Type: flightStatsByDateType,
		Args: graphql.FieldConfigArgument{
			"origin":      airportCodeArgument,
			"destination": airportCodeArgument,
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			origin, _ := params.Args["origin"].(string)
			dest, _ := params.Args["destination"].(string)

			stats, err := st.FlightStats(
				params.Context,
				origin, dest,
				store.FlightStatsOpts{
					TimeGroup: store.GroupByDay,
				},
			)

			if err == store.ErrInvalidAirportCode {
				return nil, nil
			} else if err != nil {
				return nil, err
			}

			return stats, nil
		},
	}
}

// monthlyFlightStatsQuery defines the monthlyFlightStats GraphQL query, which
// returns flight stats grouped by day.
// The store instance is used when resolving the query.
func monthlyFlightStatsQuery(st *store.Store) *graphql.Field {
	return &graphql.Field{
		Type: flightStatsByDateType,
		Args: graphql.FieldConfigArgument{
			"origin":      airportCodeArgument,
			"destination": airportCodeArgument,
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			origin, _ := params.Args["origin"].(string)
			dest, _ := params.Args["destination"].(string)

			stats, err := st.FlightStats(
				params.Context,
				origin, dest,
				store.FlightStatsOpts{
					TimeGroup: store.GroupByMonth,
				},
			)

			if err == store.ErrInvalidAirportCode {
				return nil, nil
			} else if err != nil {
				return nil, err
			}

			return stats, nil
		},
	}
}
