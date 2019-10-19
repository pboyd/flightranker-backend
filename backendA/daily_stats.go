package main

import (
	"database/sql"
	"sort"
	"strings"

	"github.com/graphql-go/graphql"
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
