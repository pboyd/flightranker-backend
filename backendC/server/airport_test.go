package server

import (
	"testing"

	"github.com/pboyd/flightranker-backend/backendC/store"
	"github.com/stretchr/testify/assert"
)

func TestAirport(t *testing.T) {
	cases := []struct {
		query    string
		expected store.Airport
	}{
		{
			query: `{airport(code:"LAS"){code,name}}`,
			expected: store.Airport{
				Code: "LAS",
				Name: "McCarran International",
			},
		},
		{
			query: `{airport(code:"LAS"){city,state}}`,
			expected: store.Airport{
				City:  "Las Vegas",
				State: "NV",
			},
		},
		{
			query:    `{airport(code:"SIX"){code,name}}`,
			expected: store.Airport{},
		},
	}

	assert := assert.New(t)

	for _, c := range cases {
		var response map[string]store.Airport
		runTestQuery(t, c.query, &response)
		assert.Equal(c.expected, response["airport"])
	}
}

func TestAirportList(t *testing.T) {
	cases := []struct {
		query         string
		expectedCodes []string
	}{
		{
			query:         `{airportList(term:"vegas"){code}}`,
			expectedCodes: []string{"LAS"},
		},
	}

	assert := assert.New(t)

	for _, c := range cases {
		var response map[string][]store.Airport
		runTestQuery(t, c.query, &response)

		actualCodes := make([]string, len(response["airportList"]))
		for i, airport := range response["airportList"] {
			actualCodes[i] = airport.Code
		}

		assert.Equal(c.expectedCodes, actualCodes)
	}
}
