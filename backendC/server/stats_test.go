package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlightStatsByAirline(t *testing.T) {
	cases := []struct {
		query    string
		expected []flightStatsByAirlineRow
	}{
		{
			query: `{flightStatsByAirline(origin:"LAS",destination:"JFK"){airline}}`,
			expected: []flightStatsByAirlineRow{
				{Airline: "JetBlue Airways"},
				{Airline: "American Airlines Inc."},
				{Airline: "Delta Air Lines Inc."},
				{Airline: "Alaska Airlines Inc."},
			},
		},
	}

	assert := assert.New(t)

	for _, c := range cases {
		var response map[string][]flightStatsByAirlineRow
		runTestQuery(t, c.query, &response)
		assert.Equal(c.expected, response["flightStatsByAirline"])
	}
}
