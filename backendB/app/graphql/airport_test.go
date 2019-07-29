package graphql

import (
	"context"
	"testing"

	"github.com/pboyd/flightranker-backend/backendb/app"
)

func TestAirport(t *testing.T) {
	cases := []struct {
		airport  *app.Airport
		query    string
		expected string
	}{
		{
			airport:  &app.Airport{Code: "SOX", Name: "Somewhere Intl"},
			query:    `{airport(code:"SOX"){code,name}}`,
			expected: `{"airport":{"code":"SOX","name":"Somewhere Intl"}}`,
		},
		{
			airport:  &app.Airport{Code: "SOX", Name: "Somewhere Intl"},
			query:    `{airport(code:"SIX"){code,name}}`,
			expected: `{"airport":null}`,
		},
	}

	for _, c := range cases {
		p := NewProcessor(ProcessorConfig{
			AirportStore: &app.AirportStoreMock{
				AirportFn: func(ctx context.Context, code string) (*app.Airport, error) {
					if c.airport.Code != code {
						return nil, nil
					}

					return c.airport, nil
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

func TestAirportList(t *testing.T) {
	cases := []struct {
		airports []*app.Airport
		query    string
		expected string
	}{
		{
			airports: []*app.Airport{
				{Code: "BAX", Name: "Baxter Intl"},
				{Code: "BEX", Name: "Beshire Regional"},
				{Code: "BIX", Name: "Bigly Municipal"},
			},
			query:    `{airportList(term:"B"){code}}`,
			expected: `{"airportList":[{"code":"BAX"},{"code":"BEX"},{"code":"BIX"}]}`,
		},
	}

	for _, c := range cases {
		p := NewProcessor(ProcessorConfig{
			AirportStore: &app.AirportStoreMock{
				AirportSearchFn: func(ctx context.Context, code string) ([]*app.Airport, error) {
					return c.airports, nil
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
