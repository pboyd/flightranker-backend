package mysql

import (
	"context"
	"reflect"
	"testing"

	"github.com/pboyd/flightranker-backend/backendb/app"
	"github.com/pboyd/flightranker-backend/backendtest"
)

func TestAirport(t *testing.T) {
	cases := []struct {
		code     string
		expected *app.Airport
	}{
		{
			code: "DEN",
			expected: &app.Airport{
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

	store := NewStoreFromDB(backendtest.ConnectMySQL(t))

	for _, c := range cases {
		actual, err := store.Airport(context.Background(), c.code)
		if err != nil {
			t.Errorf("%s: got error %v, want nil", c.code, err)
			continue
		}

		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("%s\ngot:  %#v\nwant: %#v", c.code, actual, c.expected)
			continue
		}
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
				"CEC",
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

	store := NewStoreFromDB(backendtest.ConnectMySQL(t))

	for _, c := range cases {
		actual, err := store.AirportSearch(context.Background(), c.term)
		if err != nil {
			t.Errorf("%q: got error %v, want nil", c.term, err)
			continue
		}

		if len(actual) != len(c.expectedCodes) {
			t.Errorf("%q: got %d results, want %d", c.term, len(actual), len(c.expectedCodes))
		}

		for i := range c.expectedCodes {
			if actual[i].Code != c.expectedCodes[i] {
				t.Errorf("%q-%d: got %q, want %q", c.term, i, actual[i].Code, c.expectedCodes[i])
			}
		}
	}
}
