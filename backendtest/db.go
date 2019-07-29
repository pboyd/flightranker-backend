package backendtest

import (
	"database/sql"
	"flag"
	"testing"
)

var mysqlDSN = flag.String("mysql-dsn", "", "MySQL DSN")

func ConnectMySQL(t *testing.T) *sql.DB {
	if mysqlDSN == nil || *mysqlDSN == "" {
		t.Skipf("no value for -mysql-dsn")
	}

	db, err := sql.Open("mysql", *mysqlDSN)
	if err != nil {
		t.Fatalf("could not connect to mysql: %v", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("could not connect to mysql: %v", err)
	}

	return db
}
