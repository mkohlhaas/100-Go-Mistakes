// Package db provides a tiny database client used to demonstrate how to
// categorize tests with build tags.
//
// Unit tests (no build tag) exercise validation logic without touching a real
// database. Integration tests (build tag `integration`) exercise the actual
// insert path and must be opted into explicitly.
package db

import "errors"

// errEmptyDSN is returned when Open is called with an empty data source name.
var errEmptyDSN = errors.New("db: empty dsn")

// DB holds a (parsed) data source name. In a real program this would also
// hold the live connection handle.
type DB struct {
	DSN string
}

// Open parses and validates a data source name. It does not open a network
// connection, so it is safe to exercise in unit tests.
func Open(dsn string) (*DB, error) {
	if dsn == "" {
		return nil, errEmptyDSN
	}
	return &DB{DSN: dsn}, nil
}

// Insert persists a row. A real implementation would send `row` to the
// database; here it only validates that a connection exists so integration
// tests have a concrete call to make against a live DB.
func (d *DB) Insert(row string) error {
	if d == nil {
		return errors.New("db: nil connection")
	}
	if d.DSN == "" {
		return errors.New("db: no connection")
	}
	return nil
}
