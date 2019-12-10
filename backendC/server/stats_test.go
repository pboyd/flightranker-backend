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

func TestDailyFlightStats(t *testing.T) {
	cases := []struct {
		query            string
		expectedAirlines []string
	}{
		{
			query: `{dailyFlightStats(origin:"LAS",destination:"JFK"){airline}}`,
			expectedAirlines: []string{
				"Alaska Airlines Inc.",
				"American Airlines Inc.",
				"Delta Air Lines Inc.",
				"JetBlue Airways",
			},
		},
	}

	assert := assert.New(t)

	for _, c := range cases {
		var response map[string][]flightStatsByDate
		runTestQuery(t, c.query, &response)

		actualAirlines := []string{}
		for _, row := range response["dailyFlightStats"] {
			actualAirlines = append(actualAirlines, row.Airline)
		}
		assert.Equal(c.expectedAirlines, actualAirlines)
	}
}

func TestMonthlyFlightStats(t *testing.T) {
	cases := []struct {
		query            string
		expectedAirlines []string
	}{
		{
			query: `{monthlyFlightStats(origin:"LAS",destination:"JFK"){airline}}`,
			expectedAirlines: []string{
				"Alaska Airlines Inc.",
				"American Airlines Inc.",
				"Delta Air Lines Inc.",
				"JetBlue Airways",
			},
		},
	}

	assert := assert.New(t)

	for _, c := range cases {
		var response map[string][]flightStatsByDate
		runTestQuery(t, c.query, &response)

		actualAirlines := []string{}
		for _, row := range response["monthlyFlightStats"] {
			actualAirlines = append(actualAirlines, row.Airline)
		}
		assert.Equal(c.expectedAirlines, actualAirlines)
	}
}
