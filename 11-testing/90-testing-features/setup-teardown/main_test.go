package main

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

// db is the package-wide test database, opened once in TestMain.
var db *sql.DB

// dsn points to a file-backed SQLite database so the data seeded in TestMain
// survives across the individual tests.
const dsn = "file:setup-teardown.db?mode=rwc"

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	if err := setupDB(db); err != nil {
		panic(err)
	}

	code := m.Run()

	teardownDB(dsn)
	os.Exit(code)
}

// setupDB creates the schema and seeds data once for the whole package.
func setupDB(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id    INTEGER PRIMARY KEY AUTOINCREMENT,
		name  TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	_, err := db.Exec(`
		INSERT OR IGNORE INTO users (name, email) VALUES
			('alice', 'alice@example.com'),
			('bob',   'bob@example.com')
	`)
	return err
}

// teardownDB removes the file-backed database after all tests finish.
func teardownDB(dsn string) {
	_ = dsn
	_ = os.Remove("setup-teardown.db")
}

// createConnection opens a per-test connection that is closed when the test
// finishes, even if the test fails.
func createConnection(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open connection: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})
	return conn
}

func TestUsersCount(t *testing.T) {
	conn := createConnection(t)

	var count int
	if err := conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		t.Fatalf("query users: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 users, got %d", count)
	}
}

func TestInsertThenCleanup(t *testing.T) {
	conn := createConnection(t)

	if _, err := conn.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "tmp", "tmp@example.com"); err != nil {
		t.Fatalf("insert: %v", err)
	}

	var count int
	if err := conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		t.Fatalf("query users: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 users after insert, got %d", count)
	}

	// t.Cleanup runs in LIFO order, so this delete runs before the
	// connection is closed, leaving a clean database for the next test.
	t.Cleanup(func() {
		if _, err := conn.Exec("DELETE FROM users WHERE email = ?", "tmp@example.com"); err != nil {
			t.Logf("cleanup: %v", err)
		}
	})
}
