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
			} else if err != nil {
				return nil, err
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

			return outStats, nil
		},
	}
}

// flightStatsByDateType is the GraphQL definition of the return value from
// dailyFlightStatsQuery and monthylyFlightStatsQuery.
var flightStatsByDateType = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "flightStatsByDate",
			Fields: graphql.Fields{
				"airline": &graphql.Field{Type: graphql.String},
				"rows": &graphql.Field{Type: graphql.NewList(graphql.NewObject(
					graphql.ObjectConfig{
						Name: "flightStatsByDateRow",
						Fields: graphql.Fields{
							"date":             &graphql.Field{Type: graphql.DateTime},
							"flights":          &graphql.Field{Type: graphql.Int},
							"delays":           &graphql.Field{Type: graphql.Int},
							"onTimePercentage": &graphql.Field{Type: graphql.Float},
						},
					},
				),
				)},
			},
		},
	),
)

// flightStatsByDate is the format returned by dailyFlightStatsQuery and
// monthlyFlightStatsQuery.
type flightStatsByDate struct {
	Airline string
	Rows    []flightStatsByDateRow
}

// flightStatsByDateRow is one row for an airline in flightStatsByDate,
// representing one unit of time (either one day or one month).
type flightStatsByDateRow struct {
	Date             time.Time `json:"date"`
	Flights          int       `json:"flights"`
	Delays           int       `json:"delays"`
	OnTimePercentage float64   `json:"onTimePercentage"`
}

// convertDateFlightStats converts the return value of store.FlightStats into
// the format required by dailyFlightStats and monthlyFlightStats.
func convertDateFlightStats(statsMap store.Stats) []flightStatsByDate {
	statsSlice := make([]flightStatsByDate, 0, len(statsMap))
	for airline, statsMapRows := range statsMap {
		rows := make([]flightStatsByDateRow, len(statsMapRows))
		for i, statsMapRow := range statsMapRows {
			rows[i] = flightStatsByDateRow{
				Date:             statsMapRow.Start,
				Flights:          statsMapRow.Flights,
				Delays:           statsMapRow.Delays,
				OnTimePercentage: statsMapRow.OnTime(),
			}
		}

		statsSlice = append(statsSlice, flightStatsByDate{Airline: airline, Rows: rows})
	}

	sort.Slice(statsSlice, func(i, j int) bool {
		return statsSlice[i].Airline < statsSlice[j].Airline
	})

	return statsSlice
}

// dailyFlightStatsQuery defines the dailyFlightStats GraphQL query, which
// returns flight stats grouped by day.
// The store instance is used when resolving the query.
func dailyFlightStatsQuery(st *store.Store) *graphql.Field {
	return &graphql.Field{
		Type: flightStatsByDateType,
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
					TimeGroup: store.GroupByDay,
				},
			)

			if err == store.ErrInvalidAirportCode {
				return nil, nil
			} else if err != nil {
				return nil, err
			}

			return convertDateFlightStats(stats), nil
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
					TimeGroup: store.GroupByMonth,
				},
			)

			if err == store.ErrInvalidAirportCode {
				return nil, nil
			} else if err != nil {
				return nil, err
			}

			return convertDateFlightStats(stats), nil
		},
	}
}
