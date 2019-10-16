package graphql

import (
	"context"
	"testing"
	"time"

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

func TestDailyFlightStats(t *testing.T) {
	cases := []struct {
		stats    map[string][]*app.FlightStatsByDateRow
		query    string
		expected string
	}{
		{
			stats: map[string][]*app.FlightStatsByDateRow{
				"Delta": []*app.FlightStatsByDateRow{
					{Date: date(2019, 01, 01), Flights: 100, Delays: 10},
					{Date: date(2019, 01, 02), Flights: 10, Delays: 0},
					{Date: date(2019, 01, 03), Flights: 23, Delays: 4},
				},
			},
			query:    `{dailyFlightStats(origin:"SOX",destination:"SAX"){airline,rows{date,flights,delays,onTimePercentage}}}`,
			expected: `{"dailyFlightStats":[{"airline":"Delta","rows":[{"date":"2019-01-01T00:00:00Z","delays":10,"flights":100,"onTimePercentage":90},{"date":"2019-01-02T00:00:00Z","delays":0,"flights":10,"onTimePercentage":100},{"date":"2019-01-03T00:00:00Z","delays":4,"flights":23,"onTimePercentage":82.6086956521739}]}]}`,
		},
	}

	for _, c := range cases {
		p := NewProcessor(ProcessorConfig{
			FlightStatsStore: &app.FlightStatsStoreMock{
				DailyFlightStatsFn: func(ctx context.Context, origin, dest string) (map[string][]*app.FlightStatsByDateRow, error) {
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

func TestMonthlyFlightStats(t *testing.T) {
	cases := []struct {
		stats    map[string][]*app.FlightStatsByDateRow
		query    string
		expected string
	}{
		{
			stats: map[string][]*app.FlightStatsByDateRow{
				"Delta": []*app.FlightStatsByDateRow{
					{Date: date(2019, 01, 01), Flights: 100, Delays: 10},
					{Date: date(2019, 02, 01), Flights: 10, Delays: 0},
					{Date: date(2019, 03, 01), Flights: 23, Delays: 4},
				},
			},
			query:    `{monthlyFlightStats(origin:"SOX",destination:"SAX"){airline,rows{date,flights,delays,onTimePercentage}}}`,
			expected: `{"monthlyFlightStats":[{"airline":"Delta","rows":[{"date":"2019-01-01T00:00:00Z","delays":10,"flights":100,"onTimePercentage":90},{"date":"2019-02-01T00:00:00Z","delays":0,"flights":10,"onTimePercentage":100},{"date":"2019-03-01T00:00:00Z","delays":4,"flights":23,"onTimePercentage":82.6086956521739}]}]}`,
		},
	}

	for _, c := range cases {
		p := NewProcessor(ProcessorConfig{
			FlightStatsStore: &app.FlightStatsStoreMock{
				MonthlyFlightStatsFn: func(ctx context.Context, origin, dest string) (map[string][]*app.FlightStatsByDateRow, error) {
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

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
