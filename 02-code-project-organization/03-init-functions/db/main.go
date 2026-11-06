// Package db
package db

import (
	"database/sql"
	"log"
	"os"
)

// db is a package-level global populated by init() below.
// Because it is initialized automatically, tests cannot control when the
// connection opens or substitute a mock.
// MISTAKE: globals initialized by init() are hard to test and reason about.
var db *sql.DB

// init runs automatically when this package is imported, before main() starts.
// It reads the data source name from the environment, opens a MySQL connection,
// pings the database, and panics on any failure.
//
// MISTAKES:
//   - Hidden side effect: importing this package forces a database connection
//     even if the caller never uses db.
//   - Cannot return errors: init() has no return value, so failure can only be
//     signaled with log.Panic (or log.Fatal), which terminates the process.
//   - Implicit: there is no explicit call site, so the ordering of this setup
//     relative to other packages is non-obvious.

// comment out for main.go to work
func init() {
	dataSourceName := os.Getenv("MYSQL_DATA_SOURCE_NAME")
	d, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Panic(err)
	}
	err = d.Ping()
	if err != nil {
		log.Panic(err)
	}
	db = d
}

// CreateClient is the explicit alternative to the init() above. It opens a MySQL
// connection and pings the database, but returns the client and any error so the
// caller decides when to connect and how to handle failure.
//
// PREFERRED APPROACH:
//   - Explicit: the caller invokes it deliberately (typically from main()).
//   - Errors are returned, not panicked.
//   - Testable: you can call it with different data source names or use a mock.
func CreateClient(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
