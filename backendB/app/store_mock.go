package app

import "context"

var _ AirportStore = &AirportStoreMock{}

type AirportStoreMock struct {
	AirportFn       func(ctx context.Context, code string) (*Airport, error)
	AirportSearchFn func(ctx context.Context, term string) ([]*Airport, error)
}

func (m *AirportStoreMock) Airport(ctx context.Context, code string) (*Airport, error) {
	return m.AirportFn(ctx, code)
}

func (m *AirportStoreMock) AirportSearch(ctx context.Context, term string) ([]*Airport, error) {
	return m.AirportSearchFn(ctx, term)
}

var _ FlightStatsStore = &FlightStatsStoreMock{}

type FlightStatsStoreMock struct {
	FlightStatsByAirlineFn func(ctx context.Context, origin, destination string) ([]*FlightStats, error)
}

func (m *FlightStatsStoreMock) FlightStatsByAirline(ctx context.Context, origin, destination string) ([]*FlightStats, error) {
	return m.FlightStatsByAirlineFn(ctx, origin, destination)
}
