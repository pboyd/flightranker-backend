package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pboyd/flightranker-backend/backendb/app/graphql"
	apphttp "github.com/pboyd/flightranker-backend/backendb/app/http"
	"github.com/pboyd/flightranker-backend/backendb/app/mysql"
)

func main() {
	store, err := mysql.NewStore(mysqlConfig())
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}

	http.Handle("/", newHandler(store))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func newHandler(store *mysql.Store) http.Handler {
	processor := graphql.NewProcessor(graphql.ProcessorConfig{
		AirportStore:     store,
		FlightStatsStore: store,
	})

	return &apphttp.Handler{
		Processor:       processor,
		CORSAllowOrigin: os.Getenv("CORS_ALLOW_ORIGIN"),
	}
}

func mysqlConfig() mysql.Config {
	return mysql.Config{
		Username: os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASS"),
		Address:  os.Getenv("MYSQL_ADDRESS"),
		DBName:   os.Getenv("MYSQL_DATABASE"),
	}
}
