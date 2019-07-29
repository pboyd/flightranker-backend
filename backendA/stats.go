package main

import (
	"context"
	"database/sql"
	"sort"
	"strings"
	"time"

	"github.com/graphql-go/graphql"
)

type airlineStats struct {
	Airline          string
	TotalFlights     int
	OnTimePercentage float64
	LastFlight       time.Time
}

var airlineStatsType = graphql.NewObject(
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

func resolveFlightStatsByAirline(db *sql.DB) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		origin, _ := p.Args["origin"].(string)
		origin = strings.ToUpper(origin)

		dest, _ := p.Args["destination"].(string)
		dest = strings.ToUpper(dest)

		if !isAirportCode(origin) || !isAirportCode(dest) {
			return nil, nil
		}

		stats, err := airlineFlightInfo(p.Context, db, origin, dest)
		if err != nil {
			return nil, err
		}

		delays, err := delaysByAirline(p.Context, db, origin, dest)
		if err != nil {
			return nil, err
		}

		for code := range stats {
			stats[code].OnTimePercentage = (1.0 - float64(delays[code])/float64(stats[code].TotalFlights)) * 100
		}

		statsRows := make([]*airlineStats, len(stats))
		i := 0
		for code := range stats {
			statsRows[i] = stats[code]
			i++
		}

		sort.Slice(statsRows, func(i, j int) bool {
			return statsRows[j].OnTimePercentage < statsRows[i].OnTimePercentage
		})

		return statsRows, nil
	}
}

func airlineFlightInfo(ctx context.Context, db *sql.DB, origin, dest string) (map[string]*airlineStats, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT
			carriers.code, carriers.name, total_flights, last_flight
		FROM (
			SELECT
				carrier, count(*) AS total_flights, max(date) AS last_flight
			FROM
				flights
			WHERE
				origin=? AND
				destination=?
			GROUP BY carrier
		) AS _
		INNER JOIN carriers ON carriers.code=carrier
		`,
		origin, dest)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := map[string]*airlineStats{}

	for rows.Next() {
		var (
			airline  string
			rowStats airlineStats
		)

		err := rows.Scan(&airline, &rowStats.Airline, &rowStats.TotalFlights, &rowStats.LastFlight)
		if err != nil {
			return nil, err
		}

		stats[airline] = &rowStats
	}

	return stats, nil
}

func delaysByAirline(ctx context.Context, db *sql.DB, origin, dest string) (map[string]int, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT
			carrier, count(*)
		FROM
			flights
		WHERE
			origin=? AND
			destination=? AND
			scheduled_departure_time <= departure_time AND
			scheduled_arrival_time <= arrival_time
		GROUP BY carrier`,
		origin, dest)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byAirline := map[string]int{}

	for rows.Next() {
		var (
			airline string
			count   int
		)

		err := rows.Scan(&airline, &count)
		if err != nil {
			return nil, err
		}

		byAirline[airline] = count
	}

	return byAirline, nil
}
