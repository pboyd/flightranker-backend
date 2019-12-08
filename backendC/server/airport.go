package server

import (
	"github.com/graphql-go/graphql"
	"github.com/pboyd/flightranker-backend/backendC/store"
)

// airportType is the GraphQL definition for an airport.
var airportType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Airport",
		Fields: graphql.Fields{
			"code":      &graphql.Field{Type: graphql.String},
			"name":      &graphql.Field{Type: graphql.String},
			"city":      &graphql.Field{Type: graphql.String},
			"state":     &graphql.Field{Type: graphql.String},
			"latitude":  &graphql.Field{Type: graphql.Float},
			"longitude": &graphql.Field{Type: graphql.Float},
		},
	},
)

// airportQuery defines a GraphQL query that accepts an airport code and
// responds with information about the airport.
type airportQuery struct {
	store *store.Store
}

// newAirportQuery creates a new airportQuery instance.
func newAirportQuery(store *store.Store) *airportQuery {
	return &airportQuery{store: store}
}

// Field returns a GraphQL schema definition.
func (q *airportQuery) Field() *graphql.Field {
	return &graphql.Field{
		Type:        airportType,
		Description: "get airport by code",
		Args: graphql.FieldConfigArgument{
			"code": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "airport IATA code (e.g. LAX)",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			code, _ := params.Args["code"].(string)
			airport, err := q.store.Airport(params.Context, code)

			if err == store.ErrInvalidAirportCode {
				return nil, nil
			}

			return airport, err
		},
	}
}
