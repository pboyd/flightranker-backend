package graphql

import (
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/pboyd/flightranker-backend/backendb/app"
)

var airlineFlightStatsType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "airlineFlightStats",
		Fields: graphql.Fields{
			"airline":          &graphql.Field{Type: graphql.String},
			"totalFlights":     &graphql.Field{Type: graphql.Int},
			"onTimePercentage": &graphql.Field{Type: graphql.Float, Resolve: resolveOnTimePercentage},
			"lastFlight":       &graphql.Field{Type: graphql.DateTime},
		},
	},
)

func (p *Processor) flightStatsByAirlineQuery() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(airlineFlightStatsType),
		Args: graphql.FieldConfigArgument{
			"origin": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "airport IATA code (e.g. LAX)",
			},
			"destination": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "airport IATA code (e.g. LAX)",
			},
		},
		Resolve: instrumentResolver("flightstats_by_airline", p.resolveFlightStatsByAirlineQuery),
	}
}

func (p *Processor) resolveFlightStatsByAirlineQuery(params graphql.ResolveParams) (interface{}, error) {
	origin, _ := params.Args["origin"].(string)
	origin = strings.ToUpper(origin)

	dest, _ := params.Args["destination"].(string)
	dest = strings.ToUpper(dest)

	if !app.IsAirportCode(origin) || !app.IsAirportCode(dest) {
		return nil, nil
	}

	return p.config.FlightStatsStore.FlightStatsByAirline(params.Context, origin, dest)
}
