package store

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// FlightStatsOpts contains flags that control how the FlightStats operates.
type FlightStatsOpts struct {
	// TimeGroup specifies the duration to include in each aggregate bucket.
	TimeGroup TimeGroup
}

// TimeGroup specifies an amount of time to include in the same aggregate
// bucket.
type TimeGroup int

const (
	// GroupByAvailable aggregates all available data in one bucket
	GroupByAvailable TimeGroup = iota
	// GroupByDay aggregates results by day.
	GroupByDay
	// GroupByMonth aggregates results by month.
	GroupByMonth
)

// Stats is the return value of FlightStats
//
// Map keys are airline names, values are a StatsRow slice with one member for
// each grouped time range.
//
// When TimeGroup is GroupByAvailable each slice will only have one entry.
type Stats map[string][]StatsRow

// StatsRow contains delay information for a single aggregated time period.
type StatsRow struct {
	// Start is the day of the earliest flight in the row.
	Start time.Time

	// Start is the day of the last flight in the row.
	End time.Time

	// Flights is the number of flights that occurred in the time period.
	Flights int

	// Delays is the number of delayed flights that occurred in the time
	// period.
	Delays int
}

// OnTime returns the percentage of flights that were on time.
func (row *StatsRow) OnTime() float64 {
	if row.Flights <= 0 {
		return 0
	}

	return (1.0 - float64(row.Delays)/float64(row.Flights)) * 100
}

// FlightStats returns delay information about flights from an origin airport
// to a destination.
//
// origin and destination are IATA airport codes (e.g. "LAX", "JFK"). If origin
// or destination is invalid ErrInvalidAirportCode is returned.
//
// See FlightStatsOpts for information about opts.
func (s *Store) FlightStats(ctx context.Context, origin, destination string, opts FlightStatsOpts) (Stats, error) {
	origin = strings.ToUpper(origin)
	destination = strings.ToUpper(destination)
	if !isAirportCode(origin) || !isAirportCode(destination) {
		return Stats{}, ErrInvalidAirportCode
	}

	groupBy := []string{"carriers.name"}

	switch opts.TimeGroup {
	case GroupByAvailable:
	case GroupByDay:
		groupBy = append(groupBy, "date")
	case GroupByMonth:
		groupBy = append(groupBy, "YEAR(date)", "MONTH(date)")
	default:
		return Stats{}, fmt.Errorf("invalid TimeGroup value %d", opts.TimeGroup)
	}

	query := fmt.Sprintf(`
		SELECT
			MIN(date),
			MAX(date),
			carriers.name,
			SUM(total_flights),
			SUM(IF(delayed_flights IS NULL, 0, delayed_flights)) AS delay_flights_not_null
		FROM
			flights_day
			INNER JOIN carriers ON carrier=carriers.code
		WHERE origin=? AND destination=? GROUP BY %s`,
		strings.Join(groupBy, ", "))

	rows, err := s.db.QueryContext(ctx, query, origin, destination)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := Stats{}

	for rows.Next() {
		var (
			airline string
			row     StatsRow
		)

		err := rows.Scan(&row.Start, &row.End, &airline, &row.Flights, &row.Delays)
		if err != nil {
			return nil, err
		}

		if stats[airline] == nil {
			stats[airline] = []StatsRow{}
		}

		stats[airline] = append(stats[airline], row)
	}

	return stats, nil
}
