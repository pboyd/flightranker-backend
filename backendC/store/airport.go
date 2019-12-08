package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Airport contains information about a flight destination or origin.
type Airport struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Airport looks up a single Airport by its code.
//
// If the airport is not found a nil Airport is returned.
//
// If the airport code is invalid ErrInvalidAirportCode is returned.
func (s *Store) Airport(ctx context.Context, code string) (*Airport, error) {
	code = strings.ToUpper(code)
	if !isAirportCode(code) {
		return nil, ErrInvalidAirportCode
	}

	row := s.db.QueryRow(`
		SELECT
			code, name, city, state, lat, lng
		FROM
			airports
		WHERE
			is_active=1 AND
			code=?
	`, code)

	var a Airport
	err := row.Scan(&a.Code, &a.Name, &a.City, &a.State, &a.Latitude, &a.Longitude)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching airport: %w", err)
	}

	return &a, nil
}

// AirportSearch finds airports with a name, city or code that contains the
// term.
//
// If nothing matches an empty list is returned.
func (s *Store) AirportSearch(ctx context.Context, term string) ([]*Airport, error) {
	termLike := fmt.Sprintf("%%%s%%", term)

	rows, err := s.db.Query(`
		SELECT
			code, name, city, state, lat, lng
		FROM
			airports
		WHERE
			is_active=1 AND (
				name LIKE ? OR
				city LIKE ? OR
				code LIKE ?
			)
	`, termLike, termLike, termLike)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []*Airport{}
	for rows.Next() {
		var a Airport
		err := rows.Scan(&a.Code, &a.Name, &a.City, &a.State, &a.Latitude, &a.Longitude)
		if err != nil {
			return nil, err
		}

		results = append(results, &a)
	}

	return results, nil
}

func isAirportCode(code string) bool {
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
