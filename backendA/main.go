package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	db, err := connectMySQL()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", graphqlHandler(db))
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func connectMySQL() (*sql.DB, error) {
	dsn := (&mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("MYSQL_ADDRESS"),
		DBName: os.Getenv("MYSQL_DATABASE"),

		AllowNativePasswords: true,
		ParseTime:            true,
	}).FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func graphqlHandler(db *sql.DB) http.HandlerFunc {
	schema, err := makeSchema(db)
	if err != nil {
		log.Fatalf("schema error: %v", err)
	}

	allowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")

	return func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("q"),
			Context:       r.Context(),
		})

		if allowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		}

		enc := json.NewEncoder(w)

		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(result.Errors)
			return
		}

		enc.Encode(result.Data)
	}
}

func makeSchema(db *sql.DB) (graphql.Schema, error) {
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

	dailyFlightStats := &graphql.Field{
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
		Resolve: resolveDailyFlightStats(db),
	}

	monthlyFlightStats := &graphql.Field{
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
