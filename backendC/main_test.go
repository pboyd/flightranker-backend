package main

import (
	"flag"
	"testing"

	"github.com/pboyd/flightranker-backend/backendC/server"
	"github.com/pboyd/flightranker-backend/backendC/store"
	"github.com/pboyd/flightranker-backend/backendtest"
)

var update = flag.Bool("update", false, "update golden files")

func TestStandardQueries(t *testing.T) {
	store, err := store.New(nil)
	if err != nil {
		t.Fatalf("error initializing store: %v", err)
	}

	runner := &backendtest.Runner{
		FixturePath: "../testfiles/golden",
		Update:      *update,
		Handler:     server.Handler(store),
	}

	runner.RunQuerySet(t, backendtest.StandardTestQueries)
}
