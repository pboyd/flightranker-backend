package mysql

import (
	"context"
	"sort"

	"github.com/pboyd/flightranker-backend/backendb/app"
)

func (s *Store) FlightStatsByAirline(ctx context.Context, origin, dest string) ([]*app.FlightStats, error) {
	rows, err := s.db.QueryContext(ctx,
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

	stats := []*app.FlightStats{}

	for rows.Next() {
		var row app.FlightStats

		err := rows.Scan(&row.Airline, &row.TotalFlights, &row.TotalDelays, &row.LastFlight)
		if err != nil {
			return nil, err
		}

		stats = append(stats, &row)
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[j].OnTimePercentage() < stats[i].OnTimePercentage()
	})

	return stats, nil
}
