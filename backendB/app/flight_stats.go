package app

import (
	"context"
	"time"
)

type FlightStatsStore interface {
	FlightStatsByAirline(ctx context.Context, origin, destination string) ([]*FlightStats, error)
}

type FlightStats struct {
	Airline      string
	TotalFlights int
	TotalDelays  int
	LastFlight   time.Time
}

func (fs *FlightStats) OnTimePercentage() float64 {
	if fs.TotalFlights <= 0 {
		return 0
	}

	return (1.0 - float64(fs.TotalDelays)/float64(fs.TotalFlights)) * 100
}
