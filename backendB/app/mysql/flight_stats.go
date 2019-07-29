package mysql

import (
	"context"
	"sort"

	"github.com/pboyd/flights/backend/backendb/app"
)

func (s *Store) FlightStatsByAirline(ctx context.Context, origin, dest string) ([]*app.FlightStats, error) {
	stats, err := s.airlineFlightInfo(ctx, origin, dest)
	if err != nil {
		return nil, err
	}

	delays, err := s.delaysByAirline(ctx, origin, dest)
	if err != nil {
		return nil, err
	}

	for code := range stats {
		stats[code].TotalDelays = delays[code]
	}

	statsRows := make([]*app.FlightStats, len(stats))
	i := 0
	for code := range stats {
		statsRows[i] = stats[code]
		i++
	}

	sort.Slice(statsRows, func(i, j int) bool {
		return statsRows[j].OnTimePercentage() < statsRows[i].OnTimePercentage()
	})

	return statsRows, nil
}

func (s *Store) airlineFlightInfo(ctx context.Context, origin, dest string) (map[string]*app.FlightStats, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT
			carriers.code, carriers.name, total_flights, last_flight
		FROM (
			SELECT
				carrier, count(*) AS total_flights, max(date) AS last_flight
			FROM
				flights
			WHERE
				origin=? AND
				destination=?
			GROUP BY carrier
		) AS _
		INNER JOIN carriers ON carriers.code=carrier
		`,
		origin, dest)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := map[string]*app.FlightStats{}

	for rows.Next() {
		var (
			airline  string
			rowStats app.FlightStats
		)

		err := rows.Scan(&airline, &rowStats.Airline, &rowStats.TotalFlights, &rowStats.LastFlight)
		if err != nil {
			return nil, err
		}

		stats[airline] = &rowStats
	}

	return stats, nil
}

func (s *Store) delaysByAirline(ctx context.Context, origin, dest string) (map[string]int, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT
			carrier, count(*)
		FROM
			flights
		WHERE
			origin=? AND
			destination=? AND
			scheduled_departure_time <= departure_time AND
			scheduled_arrival_time <= arrival_time
		GROUP BY carrier`,
		origin, dest)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byAirline := map[string]int{}

	for rows.Next() {
		var (
			airline string
			count   int
		)

		err := rows.Scan(&airline, &count)
		if err != nil {
			return nil, err
		}

		byAirline[airline] = count
	}

	return byAirline, nil
}
