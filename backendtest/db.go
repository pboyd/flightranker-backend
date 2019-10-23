package backendtest

import (
	"database/sql"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func ConnectMySQL(t *testing.T) *sql.DB {
	dsn := dsnFromEnv()
	if dsn == "" {
		t.Skipf("no value for MYSQL_ADDRESS and/or MYSQL_DATABASE")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("could not connect to mysql: %v", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("could not connect to mysql: %v", err)
	}

	return db
}

// dsnFromEnv reads database connection info from environment variables and
// builds a DSN string..
func dsnFromEnv() string {
	config := &mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("MYSQL_ADDRESS"),
		DBName: os.Getenv("MYSQL_DATABASE"),

		AllowNativePasswords: true,
		ParseTime:            true,
	}

	if config.Addr == "" || config.DBName == "" {
		// Apparently the variables aren't set.
		return ""
	}

	return config.FormatDSN()
}
