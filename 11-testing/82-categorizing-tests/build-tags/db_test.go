//go:build integration
// +build integration

package db

import (
	"os"
	"testing"
)

// TestInsert1 exercises the insert path when DB_DSN is set. It is skipped
// otherwise so it never runs against a shared database by accident.
func TestInsert1(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		t.Skip("DB_DSN not set; skipping integration test")
	}

	db, err := Open(dsn)
	if err != nil {
		t.Fatalf("could not open db: %v", err)
	}

	if err := db.Insert("row-1"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
}

// TestInsert2 is gated behind the INTEGRATION environment variable as a
// second safety net: even when the file is compiled with -tags integration,
// the test only runs when INTEGRATION=true.
func TestInsert2(t *testing.T) {
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("skipping integration test")
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		t.Skip("DB_DSN not set; skipping integration test")
	}

	db, err := Open(dsn)
	if err != nil {
		t.Fatalf("could not open db: %v", err)
	}

	if err := db.Insert("row-2"); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
}
