package store

import (
	"context"
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
	// Start is the beginning of the time period represented by the row.
	Start time.Time

	// End is the end of the time period represented by the row.
	End time.Time

	// Flights is the number of flights that occurred in the time period.
	Flights int

	// Delays is the number of delayed flights that occurred in the time
	// period.
	Delays int
}

// FlightStats returns delay information about flights from an origin airport
// to a destination.
//
// origin and destination are IATA airport codes (e.g. "LAX", "JFK").
//
// See FlightStatsOpts for information about opts.
func (s *Store) FlightStats(ctx context.Context, origin, destination string, opts FlightStatsOpts) (*Stats, error) {
	return nil, nil
}
