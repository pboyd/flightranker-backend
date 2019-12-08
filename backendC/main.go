package main

import (
	"log"

	"github.com/pboyd/flightranker-backend/backendC/server"
	"github.com/pboyd/flightranker-backend/backendC/store"
)

func main() {
	store, err := store.New(nil)
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	log.Fatal(server.Run(store))
}
