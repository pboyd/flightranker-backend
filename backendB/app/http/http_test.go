package http

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/pboyd/flights/backend/backendb/app"
	"github.com/pboyd/flights/backend/backendb/app/graphql"
)

func TestHandler(t *testing.T) {
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
	}

	for _, c := range cases {
		p := graphql.NewProcessor(graphql.ProcessorConfig{
			AirportStore: &app.AirportStoreMock{
				AirportFn: func(ctx context.Context, code string) (*app.Airport, error) {
					if c.airport.Code != code {
						return nil, nil
					}

					return c.airport, nil
				},
			},
		})
		h := &Handler{Processor: p}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", fmt.Sprintf("/?q=%s", c.query), nil)
		h.ServeHTTP(w, r)

		actual := w.Body.String()

		if actual != c.expected {
			t.Errorf("\ngot:  %s\nwant: %s", actual, c.expected)
		}
	}
}
