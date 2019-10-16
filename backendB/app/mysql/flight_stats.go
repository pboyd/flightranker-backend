package mysql

import (
	"context"
	"sort"
	"time"

	"github.com/pboyd/flightranker-backend/backendb/app"
)

func (s *Store) FlightStatsByAirline(ctx context.Context, origin, dest string) ([]*app.FlightStats, error) {
	stats, err := s.airlineFlightInfo(ctx, origin, dest)
	if err != nil {
		return nil, err
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[j].OnTimePercentage() < stats[i].OnTimePercentage()
	})

	return stats, nil
}

func (s *Store) airlineFlightInfo(ctx context.Context, origin, dest string) ([]*app.FlightStats, error) {
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

	return stats, nil
}

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
