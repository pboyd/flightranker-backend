package main

import (
	"database/sql"
	"time"

	"github.com/graphql-go/graphql"
)

func resolveMonthlyFlightStats(db *sql.DB) graphql.FieldResolveFn {
	return graphQLMetrics("monthly_flight_stats",
		func(p graphql.ResolveParams) (interface{}, error) {
			origin := getAirportCodeParam(p, "origin")
			dest := getAirportCodeParam(p, "destination")
			if origin == "" || dest == "" {
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
