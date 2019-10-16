package main

import (
	"database/sql"
	"sort"
	"strings"
	"time"

	"github.com/graphql-go/graphql"
)

type airlineStats struct {
	Airline          string
	TotalFlights     int
	OnTimePercentage float64
	LastFlight       time.Time
}

var airlineStatsType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "airlineFlightStats",
		Fields: graphql.Fields{
			"airline":          &graphql.Field{Type: graphql.String},
			"totalFlights":     &graphql.Field{Type: graphql.Int},
			"onTimePercentage": &graphql.Field{Type: graphql.Float},
			"lastFlight":       &graphql.Field{Type: graphql.DateTime},
		},
	},
)

func resolveFlightStatsByAirline(db *sql.DB) graphql.FieldResolveFn {
	return graphQLMetrics("flightstats_by_airline",
		func(p graphql.ResolveParams) (interface{}, error) {
			origin, _ := p.Args["origin"].(string)
			origin = strings.ToUpper(origin)

			dest, _ := p.Args["destination"].(string)
			dest = strings.ToUpper(dest)

			if !isAirportCode(origin) || !isAirportCode(dest) {
				return nil, nil
			}

			rows, err := db.QueryContext(p.Context,
				`SELECT
					carriers.name AS carrier_name, total_flights, delays_flights, last_flight
				FROM
					(
						SELECT
							carrier AS carrier_code,
							SUM(total_flights) AS total_flights,
							SUM(delayed_flights) AS delays_flights,
							MAX(date) AS last_flight
						FROM
							flights_day
						WHERE origin=? AND destination=?
						GROUP BY carrier_code
					) AS stats
				INNER JOIN carriers ON carrier_code=carriers.code
				`,
				origin, dest)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			stats := []*airlineStats{}

			for rows.Next() {
				var (
					row            airlineStats
					delayedFlights int
				)
				err := rows.Scan(&row.Airline, &row.TotalFlights, &delayedFlights, &row.LastFlight)
				if err != nil {
					return nil, err
				}

				row.OnTimePercentage = (1.0 - float64(delayedFlights)/float64(row.TotalFlights)) * 100

				stats = append(stats, &row)
			}

			sort.Slice(stats, func(i, j int) bool {
				return stats[j].OnTimePercentage < stats[i].OnTimePercentage
			})

			return stats, nil
		},
	)
}

type flightStatsByDate struct {
	Airline string
	Rows    []*flightStatsByDateRow
}

type flightStatsByDateRow struct {
	Date             time.Time
	Flights          int
	Delays           int
	OnTimePercentage float64
}

var gqlFlightStatsByDateRow = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "flightStatsByDateRow",
		Fields: graphql.Fields{
			"date":             &graphql.Field{Type: graphql.DateTime},
			"flights":          &graphql.Field{Type: graphql.Int},
			"delays":           &graphql.Field{Type: graphql.Int},
			"onTimePercentage": &graphql.Field{Type: graphql.Float},
		},
	},
)

var gqlFlightStatsByDate = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "flightStatsByDate",
			Fields: graphql.Fields{
				"airline": &graphql.Field{Type: graphql.String},
				"rows":    &graphql.Field{Type: graphql.NewList(gqlFlightStatsByDateRow)},
			},
		},
	),
)

func resolveDailyFlightStats(db *sql.DB) graphql.FieldResolveFn {
	return graphQLMetrics("daily_flight_stats",
		func(p graphql.ResolveParams) (interface{}, error) {
			origin, _ := p.Args["origin"].(string)
			origin = strings.ToUpper(origin)

			dest, _ := p.Args["destination"].(string)
			dest = strings.ToUpper(dest)

			if !isAirportCode(origin) || !isAirportCode(dest) {
				return nil, nil
			}

			rows, err := db.QueryContext(p.Context,
				`SELECT
					date,
					carriers.name,
					total_flights,
					IF(delayed_flights IS NULL, 0, delayed_flights) AS delay_flights_not_null
				FROM
					flights_day
					INNER JOIN carriers ON carrier=carriers.code
				WHERE origin=? AND destination=?
				ORDER BY date`,
				origin, dest)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			statsMap := map[string][]*flightStatsByDateRow{}

			for rows.Next() {
				var (
					airline string
					row     flightStatsByDateRow
				)

				err := rows.Scan(&row.Date, &airline, &row.Flights, &row.Delays)
				if err != nil {
					return nil, err
				}

				row.OnTimePercentage = (1.0 - float64(row.Delays)/float64(row.Flights)) * 100

				if statsMap[airline] == nil {
					statsMap[airline] = []*flightStatsByDateRow{}
				}

				statsMap[airline] = append(statsMap[airline], &row)
			}

			stats := make([]flightStatsByDate, 0, len(statsMap))
			for airline, rows := range statsMap {
				stats = append(stats, flightStatsByDate{
					Airline: airline,
					Rows:    rows,
				})
			}

			sort.Slice(stats, func(i, j int) bool {
				return stats[i].Airline < stats[j].Airline
			})

			return stats, nil
		},
	)
}

func resolveMonthlyFlightStats(db *sql.DB) graphql.FieldResolveFn {
	return graphQLMetrics("monthly_flight_stats",
		func(p graphql.ResolveParams) (interface{}, error) {
			origin, _ := p.Args["origin"].(string)
			origin = strings.ToUpper(origin)

			dest, _ := p.Args["destination"].(string)
			dest = strings.ToUpper(dest)

			if !isAirportCode(origin) || !isAirportCode(dest) {
				return nil, nil
			}

			rows, err := db.QueryContext(p.Context,
				`SELECT
					YEAR(date) AS year,
					MONTH(date) AS month,
					carriers.name,
					SUM(total_flights),
					SUM(IF(delayed_flights IS NULL, 0, delayed_flights)) AS delay_flights_not_null
				FROM
					flights_day
					INNER JOIN carriers ON carrier=carriers.code
				WHERE origin=? AND destination=? GROUP BY year, month, carriers.name`,
				origin, dest)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			statsMap := map[string][]*flightStatsByDateRow{}

			for rows.Next() {
				var (
					airline     string
					row         flightStatsByDateRow
					year, month int
				)

				err := rows.Scan(&year, &month, &airline, &row.Flights, &row.Delays)
				if err != nil {
					return nil, err
				}

				row.Date = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
				row.OnTimePercentage = (1.0 - float64(row.Delays)/float64(row.Flights)) * 100

				if statsMap[airline] == nil {
					statsMap[airline] = []*flightStatsByDateRow{}
				}

				statsMap[airline] = append(statsMap[airline], &row)
			}

			stats := make([]flightStatsByDate, 0, len(statsMap))
			for airline, rows := range statsMap {
				stats = append(stats, flightStatsByDate{
					Airline: airline,
					Rows:    rows,
				})
			}

			sort.Slice(stats, func(i, j int) bool {
				return stats[i].Airline < stats[j].Airline
			})

			return stats, nil
		},
	)
}

/*
SELECT YEAR(date) AS year, MONTH(date) AS month, carriers.name, SUM(total_flights), SUM(IF(delayed_flights IS NULL, 0, delayed_flights)) AS delay_flights_not_null FROM flights_day INNER JOIN carriers ON carrier=carriers.code WHERE origin='RDU' AND destination='LAX' GROUP BY year, month, carriers.name;
*/
