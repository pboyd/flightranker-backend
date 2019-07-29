package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pboyd/flightranker-backend/backendb/app"
)

func (s *Store) Airport(ctx context.Context, code string) (*app.Airport, error) {
	code = strings.ToUpper(code)

	row := s.db.QueryRow(`
		SELECT
			code, name, city, state, lat, lng
		FROM
			airports
		WHERE
			is_active=1 AND
			code=?
	`, code)

	var a app.Airport
	err := row.Scan(&a.Code, &a.Name, &a.City, &a.State, &a.Latitude, &a.Longitude)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &a, nil
}

func (s *Store) AirportSearch(ctx context.Context, term string) ([]*app.Airport, error) {
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

	results := []*app.Airport{}
	for rows.Next() {
		var a app.Airport
		err := rows.Scan(&a.Code, &a.Name, &a.City, &a.State, &a.Latitude, &a.Longitude)
		if err != nil {
			return nil, err
		}

		results = append(results, &a)
	}

	return results, nil
}
