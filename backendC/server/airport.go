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

// airportCodeArgument is the graphql definition for an argument that accepts
// an airport code.
var airportCodeArgument = &graphql.ArgumentConfig{
	Type:        graphql.String,
	Description: "airport IATA code (e.g. LAX)",
}

// airportQuery defines a GraphQL query that accepts an airport code and
// responds with information about the airport.
func airportQuery(st *store.Store) *graphql.Field {
	return &graphql.Field{
		Type:        airportType,
		Description: "get airport by code",
		Args: graphql.FieldConfigArgument{
			"code": airportCodeArgument,
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			code, _ := params.Args["code"].(string)
			airport, err := st.Airport(params.Context, code)

			if err == store.ErrInvalidAirportCode {
				return nil, nil
			}

			return airport, err
		},
	}
}

// airportListQuery defines a GraphQL query that accepts an airport code and
// responds with information about the airport.
func airportListQuery(st *store.Store) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(airportType),
		Description: "search airports",
		Args: graphql.FieldConfigArgument{
			"term": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "search term",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			term, _ := params.Args["term"].(string)
			airports, err := st.AirportSearch(params.Context, term)

			if err == store.ErrInvalidTerm {
				return nil, nil
			}

			return airports, err
		},
	}
}
