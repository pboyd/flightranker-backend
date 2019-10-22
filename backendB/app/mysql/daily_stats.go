package mysql

import (
	"context"

	"github.com/pboyd/flightranker-backend/backendb/app"
)

func (s *Store) DailyFlightStats(ctx context.Context, origin, destination string) (map[string][]*app.FlightStatsByDateRow, error) {
	rows, err := s.db.QueryContext(ctx,
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
		origin, destination)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := map[string][]*app.FlightStatsByDateRow{}

	for rows.Next() {
		var (
			airline string
			row     app.FlightStatsByDateRow
		)

		err := rows.Scan(&row.Date, &airline, &row.Flights, &row.Delays)
		if err != nil {
			return nil, err
		}

		if stats[airline] == nil {
			stats[airline] = []*app.FlightStatsByDateRow{}
		}

		stats[airline] = append(stats[airline], &row)
	}

	return stats, nil
}
