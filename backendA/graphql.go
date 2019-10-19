package main

import (
	"database/sql"

	"github.com/graphql-go/graphql"
)

func makeGQLSchema(db *sql.DB) (graphql.Schema, error) {
	airportType := graphql.NewObject(
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

	airportQuery := &graphql.Field{
		Type:        airportType,
		Description: "get airport by code",
		Args: graphql.FieldConfigArgument{
			"code": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "airport IATA code (e.g. LAX)",
			},
		},
		Resolve: resolveAirportQuery(db),
	}

	airportList := &graphql.Field{
		Type:        graphql.NewList(airportType),
		Description: "search airports",
		Args: graphql.FieldConfigArgument{
			"term": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "search term",
			},
		},
		Resolve: resolveAirportList(db),
	}

	airlineStatsType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "airlineFlightStats",
			Fields: graphql.Fields{
				"airline":          &graphql.Field{Type: graphql.String},
				"totalFlights":     &graphql.Field{Type: graphql.Int},
				"onTimePercentage": &graphql.Field{Type: graphql.Float},
				"lastFlight":       &graphql.Field{Type: graphql.DateTime},
			},
		},
	)

	flightStatsByAirline := &graphql.Field{
		Type: graphql.NewList(airlineStatsType),
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
		Resolve: resolveFlightStatsByAirline(db),
	}

	flightStatsByDateRow := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "flightStatsByDateRow",
			Fields: graphql.Fields{
				"date":             &graphql.Field{Type: graphql.DateTime},
				"flights":          &graphql.Field{Type: graphql.Int},
				"delays":           &graphql.Field{Type: graphql.Int},
				"onTimePercentage": &graphql.Field{Type: graphql.Float},
			},
		},
	)

	flightStatsByDate := graphql.NewList(
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

	dailyFlightStats := &graphql.Field{
		Type: flightStatsByDate,
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
		Resolve: resolveDailyFlightStats(db),
	}

	monthlyFlightStats := &graphql.Field{
		Type: flightStatsByDate,
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
		Resolve: resolveMonthlyFlightStats(db),
	}

	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "Query",
					Fields: graphql.Fields{
						"airport":              airportQuery,
						"airportList":          airportList,
						"flightStatsByAirline": flightStatsByAirline,
						"dailyFlightStats":     dailyFlightStats,
						"monthlyFlightStats":   monthlyFlightStats,
					},
				},
			),
		},
	)
}
