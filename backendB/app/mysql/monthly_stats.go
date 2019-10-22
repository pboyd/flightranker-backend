package mysql

import (
	"context"
	"time"

	"github.com/pboyd/flightranker-backend/backendb/app"
)

func (s *Store) MonthlyFlightStats(ctx context.Context, origin, destination string) (map[string][]*app.FlightStatsByDateRow, error) {
	rows, err := s.db.QueryContext(ctx,
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
		origin, destination)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := map[string][]*app.FlightStatsByDateRow{}

	for rows.Next() {
		var (
			airline     string
			row         app.FlightStatsByDateRow
			year, month int
		)

		err := rows.Scan(&year, &month, &airline, &row.Flights, &row.Delays)
		if err != nil {
			return nil, err
		}

		row.Date = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

		if stats[airline] == nil {
			stats[airline] = []*app.FlightStatsByDateRow{}
		}

		stats[airline] = append(stats[airline], &row)
	}

	return stats, nil
}
