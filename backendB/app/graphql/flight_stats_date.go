package graphql

import (
	"sort"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/pboyd/flightranker-backend/backendb/app"
)

var flightStatsByDateRow = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "flightStatsByDateRow",
		Fields: graphql.Fields{
			"date":             &graphql.Field{Type: graphql.DateTime},
			"flights":          &graphql.Field{Type: graphql.Int},
			"delays":           &graphql.Field{Type: graphql.Int},
			"onTimePercentage": &graphql.Field{Type: graphql.Float, Resolve: resolveOnTimePercentage},
		},
	},
)

var gqlFlightStatsByDate = graphql.NewList(
	graphql.NewObject(
		graphql.ObjectConfig{
			Name: "flightStatsByDate",
			Fields: graphql.Fields{
				"airline": &graphql.Field{Type: graphql.String},
				"rows":    &graphql.Field{Type: graphql.NewList(flightStatsByDateRow)},
			},
		},
	),
)

func (p *Processor) dailyFlightStatsQuery() *graphql.Field {
	return &graphql.Field{
		Type: gqlFlightStatsByDate,
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

type flightStatsByDate struct {
	Airline string
	Rows    []*app.FlightStatsByDateRow
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

	stats := make([]flightStatsByDate, 0, len(statsMap))
	for airline, rows := range statsMap {
		stats = append(stats, flightStatsByDate{Airline: airline, Rows: rows})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Airline < stats[j].Airline
	})

	return stats, nil
}

func (p *Processor) monthlyFlightStatsQuery() *graphql.Field {
	return &graphql.Field{
		Type: gqlFlightStatsByDate,
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
		Resolve: instrumentResolver("monthly_flight_stats", p.resolveMonthlyFlightStats),
	}
}

func (p *Processor) resolveMonthlyFlightStats(params graphql.ResolveParams) (interface{}, error) {
	origin, _ := params.Args["origin"].(string)
	origin = strings.ToUpper(origin)

	dest, _ := params.Args["destination"].(string)
	dest = strings.ToUpper(dest)

	if !app.IsAirportCode(origin) || !app.IsAirportCode(dest) {
		return nil, nil
	}

	statsMap, err := p.config.FlightStatsStore.MonthlyFlightStats(params.Context, origin, dest)
	if err != nil {
		return nil, err
	}

	stats := make([]flightStatsByDate, 0, len(statsMap))
	for airline, rows := range statsMap {
		stats = append(stats, flightStatsByDate{Airline: airline, Rows: rows})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Airline < stats[j].Airline
	})

	return stats, nil
}
