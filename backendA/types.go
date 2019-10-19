package main

import (
	"time"
)

type airport struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	City      string `json:"city"`
	State     string `json:"state"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type airlineStats struct {
	Airline          string
	TotalFlights     int
	OnTimePercentage float64
	LastFlight       time.Time
}

type flightStatsByDate struct {
	Airline string
	Rows    []*flightStatsByDateRow
}

type flightStatsByDateRow struct {
	Date             time.Time
	Flights          int
	Delays           int
	OnTimePercentage float64
}
