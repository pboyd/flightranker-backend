package mysql

import (
	"context"
	"sort"
	"testing"

	"github.com/pboyd/flightranker-backend/backendb/app"
	"github.com/pboyd/flightranker-backend/backendtest"
)

func TestFlightStatsByAirline(t *testing.T) {
	cases := []struct {
		origin, dest string
		expected     []*app.FlightStats
	}{
		{
			origin: "DEN",
			dest:   "LAS",
			expected: []*app.FlightStats{
				{Airline: "Frontier Airlines Inc."},
				{Airline: "Southwest Airlines Co."},
				{Airline: "Spirit Air Lines"},
				{Airline: "United Air Lines Inc."},
			},
		},
	}

	store := NewStoreFromDB(backendtest.ConnectMySQL(t))

	for _, c := range cases {
		actual, err := store.FlightStatsByAirline(context.Background(), c.origin, c.dest)
		if err != nil {
			t.Errorf("%s-%s: got error %v, want nil", c.origin, c.dest, err)
			continue
		}

		if len(actual) != len(c.expected) {
			t.Errorf("%s-%s: got %d results, want %d", c.origin, c.dest, len(actual), len(c.expected))
		}

		sort.Slice(actual, func(i, j int) bool {
			return actual[i].Airline < actual[j].Airline
		})

		for i := range c.expected {
			if i >= len(actual) {
				t.Errorf("%s-%s-%d: missing item", c.origin, c.dest, i)
				continue
			}

			if actual[i].Airline != c.expected[i].Airline {
				t.Errorf("%s-%s-%d: got Airline %q, want %q", c.origin, c.dest, i, actual[i].Airline, c.expected[i].Airline)
			}

			if actual[i].TotalFlights <= 0 {
				t.Errorf("%s-%s-%d: got TotalFlights %d, want >0", c.origin, c.dest, i, actual[i].TotalFlights)
			}

			if actual[i].TotalDelays <= 0 {
				t.Errorf("%s-%s-%d: got TotalDelays %d, want >0", c.origin, c.dest, i, actual[i].TotalDelays)
			}
		}
	}
}

func TestDailyFlightStats(t *testing.T) {
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

	store := NewStoreFromDB(backendtest.ConnectMySQL(t))

	for _, c := range cases {
		actual, err := store.DailyFlightStats(context.Background(), c.origin, c.dest)
		if err != nil {
			t.Errorf("%s-%s: got error %v, want nil", c.origin, c.dest, err)
			continue
		}

		if len(actual) != len(c.expectedAirlines) {
			t.Errorf("%s-%s: got %d airlines, want %d", c.origin, c.dest, len(actual), len(c.expectedAirlines))
		}

		for _, airline := range c.expectedAirlines {
			series := actual[airline]
			if series == nil {
				t.Errorf("%s-%s-%q: missing airline", c.origin, c.dest, airline)
				continue
			}

			if len(series) == 0 {
				t.Errorf("%s-%s-%q: empty series", c.origin, c.dest, airline)
				continue
			}

			for _, dayStats := range series {
				if dayStats.Date.IsZero() {
					t.Errorf("%s-%s-%q: contains zero date", c.origin, c.dest, airline)
				}
			}
		}
	}
}
