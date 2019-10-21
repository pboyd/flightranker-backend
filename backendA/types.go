package main

import (
	"sort"
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

type flightStatsByDateSlice []flightStatsByDate

func (s flightStatsByDateSlice) Sort() {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Airline < s[j].Airline
	})
}

// newFlightStatsByDateSlice converts a map of airline names and flight stats to a flightStatsByDateSlice
func newFlightStatsByDateSlice(m map[string][]*flightStatsByDateRow) flightStatsByDateSlice {
	stats := make(flightStatsByDateSlice, 0, len(m))
	for airline, rows := range m {
		stats = append(stats, flightStatsByDate{
			Airline: airline,
			Rows:    rows,
		})
	}

	return stats
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
