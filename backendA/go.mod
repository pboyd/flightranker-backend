module github.com/pboyd/flightranker-backend/backendA

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/graphql-go/graphql v0.7.8
	github.com/pboyd/flightranker-backend/backendtest v0.0.0
	google.golang.org/appengine v1.6.1 // indirect
)

replace github.com/pboyd/flightranker-backend/backendtest => ../backendtest
