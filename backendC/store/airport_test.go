package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAirport(t *testing.T) {
	cases := []struct {
		code     string
		expected *Airport
	}{
		{
			code: "DEN",
			expected: &Airport{
				Code:      "DEN",
				Name:      "Denver Intl",
				City:      "Denver",
				State:     "CO",
				Latitude:  39.85840806,
				Longitude: -104.66700190,
			},
		},
		{
			code:     "XYZ",
			expected: nil,
		},
	}

	store := testStore(t)

	assert := assert.New(t)

	for _, c := range cases {
		actual, err := store.Airport(context.Background(), c.code)
		if !assert.NoError(err) {
			continue
		}

		assert.Equal(c.expected, actual)
	}
}

func TestAirportSearch(t *testing.T) {
	cases := []struct {
		term          string
		expectedCodes []string
	}{
		{
			term: "jack",
			expectedCodes: []string{
				"JAC",
				"JAN",
				"JAX",
				"OAJ",
			},
		},
		{
			term:          "XYZ",
			expectedCodes: []string{},
		},
	}

	store := testStore(t)
	assert := assert.New(t)

	for _, c := range cases {
		actual, err := store.AirportSearch(context.Background(), c.term)
		if !assert.NoError(err) {
			continue
		}

		actualCodes := make([]string, len(actual))
		for i, airport := range actual {
			actualCodes[i] = airport.Code
		}

		assert.Equal(c.expectedCodes, actualCodes)
	}
}
