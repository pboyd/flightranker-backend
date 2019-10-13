package app

import (
	"context"
	"time"
)

type FlightStatsStore interface {
	FlightStatsByAirline(ctx context.Context, origin, destination string) ([]*FlightStats, error)
	DailyFlightStats(ctx context.Context, origin, destination string) (map[string][]*FlightStatsDay, error)
}

type FlightStats struct {
	Airline      string
	TotalFlights int
	TotalDelays  int
	LastFlight   time.Time
}

func (fs *FlightStats) OnTimePercentage() float64 {
	return calcOnTimePercentage(fs.TotalFlights, fs.TotalDelays)
}

type FlightStatsDay struct {
	Date    time.Time
	Flights int
	Delays  int
}

func (fs *FlightStatsDay) OnTimePercentage() float64 {
	return calcOnTimePercentage(fs.Flights, fs.Delays)
}

func calcOnTimePercentage(total, delayed int) float64 {
	if total <= 0 {
		return 0
	}

	return (1.0 - float64(delayed)/float64(total)) * 100
}

type OnTimeStat interface {
	OnTimePercentage() float64
}
