package app

import (
	"context"
	"strings"
	"unicode"
)

type AirportStore interface {
	Airport(ctx context.Context, code string) (*Airport, error)
	AirportSearch(ctx context.Context, term string) ([]*Airport, error)
}

type Airport struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// IsAirportCode returns true if the airport is a syntactically valid IATA
// airport code (3 uppercase Latin letters).
func IsAirportCode(code string) bool {
	if len(code) != 3 {
		return false
	}

	for _, r := range code {
		if r < 'A' || r > 'Z' {
			return false
		}
	}

	return true
}

// IsValidAirportSearchTerm returns true if the term is valid. Specifically,
// containing only digits 0 to 9, Latin letters, dashes or spaces.
func IsValidAirportSearchTerm(term string) bool {
	if term == "" {
		return false
	}

	badStart := strings.TrimLeftFunc(term, func(r rune) bool {
		if r > unicode.MaxLatin1 {
			return false
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return true
		}

		if r == '-' || r == ' ' {
			return true
		}

		return false
	})

	return len(badStart) == 0
}
