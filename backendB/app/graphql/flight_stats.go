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

func resolveOnTimePercentage(params graphql.ResolveParams) (interface{}, error) {
	stats, ok := params.Source.(app.OnTimeStat)
	if !ok {
		return 0, nil
	}

	return stats.OnTimePercentage(), nil
}

var dailyFlightStatsRow = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "dailyFlightStatsDay",
		Fields: graphql.Fields{
			"date":             &graphql.Field{Type: graphql.DateTime},
			"flights":          &graphql.Field{Type: graphql.Int},
			"delays":           &graphql.Field{Type: graphql.Int},
			"onTimePercentage": &graphql.Field{Type: graphql.Float, Resolve: resolveOnTimePercentage},
		},
	},
)

var dailyFlightStats = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "dailyFlightStats",
			Fields: graphql.Fields{
				"airline": &graphql.Field{Type: graphql.String},
				"days":    &graphql.Field{Type: graphql.NewList(dailyFlightStatsRow)},
			},
		},
	),
)

func (p *Processor) dailyFlightStatsQuery() *graphql.Field {
	return &graphql.Field{
		Type: dailyFlightStats,
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
		Resolve: instrumentResolver("daily_flight_stats", p.resolveDailyFlightStats),
	}
}

type dailyStatsAirline struct {
	Airline string
	Days    []*app.FlightStatsDay
}

func (p *Processor) resolveDailyFlightStats(params graphql.ResolveParams) (interface{}, error) {
	origin, _ := params.Args["origin"].(string)
	origin = strings.ToUpper(origin)

	dest, _ := params.Args["destination"].(string)
	dest = strings.ToUpper(dest)

	if !app.IsAirportCode(origin) || !app.IsAirportCode(dest) {
		return nil, nil
	}

	statsMap, err := p.config.FlightStatsStore.DailyFlightStats(params.Context, origin, dest)
	if err != nil {
		return nil, err
	}

	stats := make([]dailyStatsAirline, 0, len(statsMap))
	for airline, days := range statsMap {
		stats = append(stats, dailyStatsAirline{Airline: airline, Days: days})
	}

	return stats, nil
}
