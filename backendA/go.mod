module github.com/pboyd/flightranker-backend/backendA

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/graphql-go/graphql v0.7.8
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pboyd/flightranker-backend/backendtest v0.0.0
	github.com/prometheus/client_golang v1.1.0
	google.golang.org/appengine v1.6.1 // indirect
)

replace github.com/pboyd/flightranker-backend/backendtest => ../backendtest
