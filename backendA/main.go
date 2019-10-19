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
	schema, err := makeGQLSchema(db)
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
