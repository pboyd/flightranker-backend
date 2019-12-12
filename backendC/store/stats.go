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

// Stats is the return value of FlightStats. Each entry in the slice contains
// data for one airline.
type Stats []AirlineStats

// AirlineStats contains data for one airline.
//
// When TimeGroup is GroupByAvailable there will only be one row. In all other
// cases, there will be one row per time period.
type AirlineStats struct {
	Airline string
	Rows    []StatsRow
}

// StatsRow contains delay information for a single aggregated time period.
type StatsRow struct {
	// Start is the day of the earliest flight in the row.
	Start time.Time `json:"date"`

	// Start is the day of the last flight in the row.
	End time.Time `json:"end_date"`

	// Flights is the number of flights that occurred in the time period.
	Flights int `json:"flights"`

	// Delays is the number of delayed flights that occurred in the time
	// period.
	Delays int `json:"delays"`
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
		WHERE origin=? AND destination=?
		GROUP BY %s
		ORDER BY carriers.name`,
		strings.Join(groupBy, ", "))

	rows, err := s.db.QueryContext(ctx, query, origin, destination)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := Stats{}
	var currentAirline *AirlineStats

	for rows.Next() {
		var (
			airline string
			row     StatsRow
		)

		err := rows.Scan(&row.Start, &row.End, &airline, &row.Flights, &row.Delays)
		if err != nil {
			return nil, err
		}

		if currentAirline == nil {
			currentAirline = &AirlineStats{
				Airline: airline,
				Rows:    []StatsRow{},
			}
		} else if airline != currentAirline.Airline {
			stats = append(stats, *currentAirline)
			currentAirline = &AirlineStats{
				Airline: airline,
				Rows:    []StatsRow{},
			}
		}

		currentAirline.Rows = append(currentAirline.Rows, row)
	}

	stats = append(stats, *currentAirline)

	return stats, nil
}
