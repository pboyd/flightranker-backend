module github.com/pboyd/flights/backend/backendb

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/graphql-go/graphql v0.7.8
	github.com/pboyd/flights/backend/backendtest v0.0.0
)

replace github.com/pboyd/flights/backend/backendtest => ../backendtest
