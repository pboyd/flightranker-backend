package main

import (
	"flag"
	"testing"

	"github.com/pboyd/flights/backend/backendb/app/mysql"
	"github.com/pboyd/flights/backend/backendtest"
)

var update = flag.Bool("update", false, "update golden files")

func TestStandardQueries(t *testing.T) {
	store := mysql.NewStoreFromDB(backendtest.ConnectMySQL(t))
	runner := &backendtest.Runner{
		FixturePath: "../testfiles/golden",
		Update:      *update,
		Handler:     newHandler(store),
	}

	runner.RunQuerySet(t, backendtest.StandardTestQueries)
}
