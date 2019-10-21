package main

import (
	"database/sql"

	"github.com/graphql-go/graphql"
)

func resolveDailyFlightStats(db *sql.DB) graphql.FieldResolveFn {
	return graphQLMetrics("daily_flight_stats",
		func(p graphql.ResolveParams) (interface{}, error) {
			origin := getAirportCodeParam(p, "origin")
			dest := getAirportCodeParam(p, "destination")
			if origin == "" || dest == "" {
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

				row.OnTimePercentage = calculateOnTimePercentage(row.Delays, row.Flights)

				if statsMap[airline] == nil {
					statsMap[airline] = []*flightStatsByDateRow{}
				}

				statsMap[airline] = append(statsMap[airline], &row)
			}

			stats := newFlightStatsByDateSlice(statsMap)
			stats.Sort()

			return stats, nil
		},
	)
}
