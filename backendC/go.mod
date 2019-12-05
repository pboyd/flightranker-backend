module github.com/pboyd/flightranker-backend/backendC

go 1.13

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/pboyd/flightranker-backend/backendtest v0.0.0
	github.com/stretchr/testify v1.4.0
)

replace github.com/pboyd/flightranker-backend/backendtest => ../backendtest
