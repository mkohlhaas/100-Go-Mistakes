package db

import (
	"errors"
	"testing"
)

func TestOpenEmptyDSN(t *testing.T) {
	_, err := Open("")
	if err == nil {
		t.Fatal("expected an error for an empty DSN")
	}
	if !errors.Is(err, errEmptyDSN) {
		t.Errorf("expected errEmptyDSN, got %v", err)
	}
}

func TestOpenValidDSN(t *testing.T) {
	db, err := Open("host=localhost port=5432 dbname=mydb sslmode=disable")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if db.DSN != "host=localhost port=5432 dbname=mydb sslmode=disable" {
		t.Errorf("expected DSN to be preserved, got %q", db.DSN)
	}
}

func TestInsertNilDB(t *testing.T) {
	var db *DB
	if err := db.Insert("row"); err == nil {
		t.Fatal("expected an error for a nil DB")
	}
}

func TestInsertNoConnection(t *testing.T) {
	db := &DB{}
	if err := db.Insert("row"); err == nil {
		t.Fatal("expected an error for a DB with no connection")
	}
}
