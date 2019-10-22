package graphql

import (
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/pboyd/flightranker-backend/backendb/app"
)

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

func (p *Processor) airportQuery() *graphql.Field {
	return &graphql.Field{
		Type:        airportType,
		Description: "get airport by code",
		Args: graphql.FieldConfigArgument{
			"code": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "airport IATA code (e.g. LAX)",
			},
		},
		Resolve: instrumentResolver("airport", p.resolveAirportQuery),
	}
}

func (p *Processor) resolveAirportQuery(params graphql.ResolveParams) (interface{}, error) {
	code := p.getAirportCodeParam(params, "code")
	if code == "" {
		return nil, nil
	}

	return p.config.AirportStore.Airport(params.Context, code)
}

func (p *Processor) airportListQuery() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(airportType),
		Description: "search airports",
		Args: graphql.FieldConfigArgument{
			"term": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "search term",
			},
		},
		Resolve: instrumentResolver("airport_list", p.resolveAirportList),
	}
}

func (p *Processor) resolveAirportList(params graphql.ResolveParams) (interface{}, error) {
	term, _ := params.Args["term"].(string)
	if !app.IsValidAirportSearchTerm(term) {
		return nil, nil
	}

	return p.config.AirportStore.AirportSearch(params.Context, term)
}

func (p *Processor) getAirportCodeParam(params graphql.ResolveParams, key string) string {
	code, _ := params.Args[key].(string)
	code = strings.ToUpper(code)
	if !app.IsAirportCode(code) {
		return ""
	}

	return code
}
