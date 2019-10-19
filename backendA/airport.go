package main

import (
	"database/sql"
	"fmt"
	"strings"
	"unicode"

	"github.com/graphql-go/graphql"
)

func resolveAirportQuery(db *sql.DB) graphql.FieldResolveFn {
	return graphQLMetrics("airport",
		func(p graphql.ResolveParams) (interface{}, error) {
			code, _ := p.Args["code"].(string)
			code = strings.ToUpper(code)
			if !isAirportCode(code) {
				return nil, nil
			}

			row := db.QueryRow(`
			SELECT
				code, name, city, state, lat, lng
			FROM
				airports
			WHERE
				is_active=1 AND
				code=?
		`, code)

			var a airport
			err := row.Scan(&a.Code, &a.Name, &a.City, &a.State, &a.Latitude, &a.Longitude)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, nil
				}
				return nil, err
			}

			return &a, nil
		},
	)
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

func resolveAirportList(db *sql.DB) graphql.FieldResolveFn {
	return graphQLMetrics("airport_list",
		func(p graphql.ResolveParams) (interface{}, error) {
			term, _ := p.Args["term"].(string)
			if !checkAirportSearchTerm(term) {
				return nil, nil
			}

			termLike := fmt.Sprintf("%%%s%%", term)

			rows, err := db.Query(`
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

			results := []*airport{}
			for rows.Next() {
				var a airport
				err := rows.Scan(&a.Code, &a.Name, &a.City, &a.State, &a.Latitude, &a.Longitude)
				if err != nil {
					return nil, err
				}

				results = append(results, &a)
			}

			return results, nil
		},
	)
}

func checkAirportSearchTerm(term string) bool {
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
