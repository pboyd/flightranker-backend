package main

import (
	"flag"
	"testing"

	"github.com/pboyd/flightranker-backend/backendC/server"
	"github.com/pboyd/flightranker-backend/backendtest"
)

var update = flag.Bool("update", false, "update golden files")

func TestStandardQueries(t *testing.T) {
	runner := &backendtest.Runner{
		FixturePath: "../testfiles/golden",
		Update:      *update,
		Handler:     server.Handler(),
	}

	runner.RunQuerySet(t, backendtest.StandardTestQueries)
}
