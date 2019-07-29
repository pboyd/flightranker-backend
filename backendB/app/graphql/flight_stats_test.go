package graphql

import (
	"context"
	"testing"

	"github.com/pboyd/flightranker-backend/backendb/app"
)

func TestFlightStatsByAirline(t *testing.T) {
	cases := []struct {
		stats    []*app.FlightStats
		query    string
		expected string
	}{
		{
			stats: []*app.FlightStats{
				{Airline: "Delta", TotalFlights: 100, TotalDelays: 10},
			},
			query:    `{flightStatsByAirline(origin:"SOX",destination:"SAX"){airline,onTimePercentage}}`,
			expected: `{"flightStatsByAirline":[{"airline":"Delta","onTimePercentage":90}]}`,
		},
	}

	for _, c := range cases {
		p := NewProcessor(ProcessorConfig{
			FlightStatsStore: &app.FlightStatsStoreMock{
				FlightStatsByAirlineFn: func(ctx context.Context, origin, dest string) ([]*app.FlightStats, error) {
					return c.stats, nil
				},
			},
		})

		actual, err := p.Do(context.Background(), c.query)
		if err != nil {
			t.Errorf("got error %v, want nil", err)
			continue
		}

		if actual != c.expected {
			t.Errorf("\ngot:  %s\nwant: %s", actual, c.expected)
		}
	}
}
