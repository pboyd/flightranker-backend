package main

import (
	"flag"
	"testing"

	"github.com/pboyd/flightranker-backend/backendtest"
)

var update = flag.Bool("update", false, "update golden files")

func TestStandardQueries(t *testing.T) {
	db := backendtest.ConnectMySQL(t)
	runner := &backendtest.Runner{
		FixturePath: "../testfiles/golden",
		Update:      *update,
		Handler:     graphqlHandler(db),
	}

	runner.RunQuerySet(t, backendtest.StandardTestQueries)
}
