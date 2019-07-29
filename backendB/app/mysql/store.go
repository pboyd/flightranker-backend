package mysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/pboyd/flights/backend/backendb/app"
)

var _ app.AirportStore = &Store{}
var _ app.FlightStatsStore = &Store{}

type Store struct {
	db *sql.DB
}

func NewStore(cfg Config) (*Store, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func NewStoreFromDB(db *sql.DB) *Store {
	return &Store{db: db}
}

type Config struct {
	Username string
	Password string
	Address  string
	DBName   string
}

func (cfg Config) DSN() string {
	return (&mysql.Config{
		User:   cfg.Username,
		Passwd: cfg.Password,
		Addr:   cfg.Address,
		DBName: cfg.DBName,

		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}).FormatDSN()
}
