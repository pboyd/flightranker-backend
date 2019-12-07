package store

import (
	"context"
	"sort"
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

	store := testStore(t)
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
		for name := range actual {
			actualAirlines = append(actualAirlines, name)
		}
		sort.Strings(actualAirlines)

		assert.Equal(c.expectedAirlines, actualAirlines)

		for _, set := range actual {
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

	store := testStore(t)
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
		for name := range actual {
			actualAirlines = append(actualAirlines, name)
		}
		sort.Strings(actualAirlines)

		assert.Equal(c.expectedAirlines, actualAirlines)

		for _, set := range actual {
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

	store := testStore(t)
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
		for name := range actual {
			actualAirlines = append(actualAirlines, name)
		}
		sort.Strings(actualAirlines)

		assert.Equal(c.expectedAirlines, actualAirlines)

		for _, set := range actual {
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
