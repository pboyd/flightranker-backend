package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlightStatsAvailable(t *testing.T) {
	cases := []struct {
		origin, dest     string
		expectedAirlines []string
	}{
		{
			origin: "DEN",
			dest:   "LAS",
			expectedAirlines: []string{
				"Frontier Airlines Inc.",
				"Southwest Airlines Co.",
				"Spirit Air Lines",
				"United Air Lines Inc.",
			},
		},
	}

	store := New()
	assert := assert.New(t)

	for _, c := range cases {
		actual, err := store.FlightStats(
			context.Background(),
			c.origin, c.dest,
			FlightStatsOpts{TimeGroup: GroupByAvailable},
		)
		if !assert.NoError(err) {
			continue
		}

		actualAirlines := make([]string, 0, len(actual))
		for _, row := range actual {
			actualAirlines = append(actualAirlines, row.Airline)
		}

		assert.Equal(c.expectedAirlines, actualAirlines)

		for _, airline := range actual {
			set := airline.Rows
			if !assert.Len(set, 1) {
				continue
			}

			assert.GreaterOrEqual(set[0].Flights, 0)
			assert.GreaterOrEqual(set[0].Delays, 0)
		}
	}
}

func TestFlightStatsDaily(t *testing.T) {
	cases := []struct {
		origin, dest     string
		expectedAirlines []string
	}{
		{
			origin: "DEN",
			dest:   "LAS",
			expectedAirlines: []string{
				"Frontier Airlines Inc.",
				"Southwest Airlines Co.",
				"Spirit Air Lines",
				"United Air Lines Inc.",
			},
		},
	}

	store := New()
	assert := assert.New(t)

	for _, c := range cases {
		actual, err := store.FlightStats(
			context.Background(),
			c.origin, c.dest,
			FlightStatsOpts{TimeGroup: GroupByDay},
		)
		if !assert.NoError(err) {
			continue
		}

		actualAirlines := make([]string, 0, len(actual))
		for _, row := range actual {
			actualAirlines = append(actualAirlines, row.Airline)
		}

		assert.Equal(c.expectedAirlines, actualAirlines)

		for _, airline := range actual {
			set := airline.Rows
			if !assert.NotEmpty(set) {
				continue
			}

			for _, dayStats := range set {
				assert.False(dayStats.Start.IsZero())
				assert.False(dayStats.End.IsZero())
			}
		}
	}
}

func TestFlightStatsMonthly(t *testing.T) {
	cases := []struct {
		origin, dest     string
		expectedAirlines []string
	}{
		{
			origin: "DEN",
			dest:   "LAS",
			expectedAirlines: []string{
				"Frontier Airlines Inc.",
				"Southwest Airlines Co.",
				"Spirit Air Lines",
				"United Air Lines Inc.",
			},
		},
	}

	store := New()
	assert := assert.New(t)

	for _, c := range cases {
		actual, err := store.FlightStats(
			context.Background(),
			c.origin, c.dest,
			FlightStatsOpts{TimeGroup: GroupByMonth},
		)
		if !assert.NoError(err) {
			continue
		}

		actualAirlines := make([]string, 0, len(actual))
		for _, airline := range actual {
			actualAirlines = append(actualAirlines, airline.Airline)
		}

		assert.Equal(c.expectedAirlines, actualAirlines)

		for _, airline := range actual {
			set := airline.Rows
			if !assert.NotEmpty(set) {
				continue
			}

			for _, monthlyStats := range set {
				assert.False(monthlyStats.Start.IsZero())
				assert.False(monthlyStats.End.IsZero())
			}
		}
	}
}
