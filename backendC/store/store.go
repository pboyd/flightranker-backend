package store

import (
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
)

// Store contains methods for retrieving flight data from the database.
type Store struct {
	db *sql.DB
}

// New creates a new Store instance based on the given config.
//
// If config is nil, DefaultConfig will be used. If New is unable to connect to
// the database an error will be returned.
func New(cfg *Config) (*Store, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	db, err := sql.Open("mysql", cfg.dsn())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

// Config contains options and connection information for the Store.
type Config struct {
	// Network address for the database (e.g. 127.0.0.1:3306)
	MysqlAddress string

	// Database name
	MysqlDatabase string

	// Username for MySQL
	MysqlUser string

	// Password for the MySQL user
	MysqlPassword string
}

// DefaultConfig creates a new config instance with values read from
// environment variables.
//
// Fields are set from the following environment variables:
//
//   - MysqlAddress: $MYSQL_ADDRESS
//   - MysqlDatabase: $MYSQL_DATABASE
//   - MysqlUser: $MYSQL_USER
//   - MysqlPassword: $MYSQL_PASS
//
func DefaultConfig() *Config {
	return &Config{
		MysqlAddress:  os.Getenv("MYSQL_ADDRESS"),
		MysqlDatabase: os.Getenv("MYSQL_DATABASE"),
		MysqlUser:     os.Getenv("MYSQL_USER"),
		MysqlPassword: os.Getenv("MYSQL_PASS"),
	}
}

// dsn returns a database connection string based on the config values.
func (cfg *Config) dsn() string {
	return (&mysql.Config{
		User:   cfg.MysqlUser,
		Passwd: cfg.MysqlPassword,
		Addr:   cfg.MysqlAddress,
		DBName: cfg.MysqlDatabase,

		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}).FormatDSN()
}
