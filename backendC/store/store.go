package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

// ErrInvalidAirportCode is returned when an airport code is invalid. To be
// valid an airport code must contain exactly three letters.
var ErrInvalidAirportCode = errors.New("invalid airport code")

// ErrInvalidTerm is returned by AirportList when the search contains invalid
// characters. Search terms must contain only letters, numbers, dashes and
// spaces.
var ErrInvalidTerm = errors.New("invalid search term")

// Store contains methods for retrieving flight data from the database.
type Store struct {
	db *sql.DB
}

// New creates a new Store instance using MySQL connection information from the
// following environment variables:
//
//   - $MYSQL_ADDRESS - Network address for the database (e.g. 127.0.0.1:3306)
//   - $MYSQL_DATABASE - Database name
//   - $MYSQL_USER - Username for MySQL
//   - $MYSQL_PASS - Password for the MySQL user
//
// If New is unable to connect to the database it will panic.
func New() *Store {
	dsn := (&mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASS"),
		Addr:   os.Getenv("MYSQL_ADDRESS"),
		DBName: os.Getenv("MYSQL_DATABASE"),

		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}).FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("unable to connect to MySQL: %v", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("unable to ping MySQL: %v", err))
	}

	return &Store{
		db: db,
	}
}
