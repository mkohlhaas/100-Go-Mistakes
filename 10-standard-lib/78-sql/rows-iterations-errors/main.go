package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
)

func get1(ctx context.Context, db *sql.DB, id string) (string, int, error) {
	// ⚠️ the error (linter recognizes it)
	rows, err := db.QueryContext(ctx, "SELECT DEP, AGE FROM EMP WHERE ID = ?", id)
	if err != nil {
		return "", 0, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v\n", err)
		}
	}()

	var (
		department string
		age        int
	)
	for rows.Next() {
		err := rows.Scan(&department, &age)
		if err != nil {
			return "", 0, err
		}
	}

	return department, age, nil
}

func get2(ctx context.Context, db *sql.DB, id string) (string, int, error) {
	rows, err := db.QueryContext(ctx,
		"SELECT DEP, AGE FROM EMP WHERE ID = ?", id)
	if err != nil {
		return "", 0, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v\n", err)
		}
	}()

	var (
		department string
		age        int
	)
	for rows.Next() {
		err := rows.Scan(&department, &age)
		if err != nil {
			return "", 0, err
		}
	}
	if err := rows.Err(); err != nil { // ✅ the fix
		return "", 0, err
	}

	return department, age, nil
}

// ---- fake driver: a self-contained in-memory "database" (no real DB needed) ----

func init() {
	sql.Register("fake", fakeDriver{})
}

type fakeDriver struct{}

func (fakeDriver) Open(string) (driver.Conn, error) { return fakeConn{}, nil }

type fakeConn struct{}

func (fakeConn) Prepare(string) (driver.Stmt, error) { return fakeStmt{}, nil }
func (fakeConn) Close() error                        { return nil }
func (fakeConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions not supported")
}

type fakeStmt struct{}

func (fakeStmt) NumInput() int { return 1 }
func (fakeStmt) Close() error  { return nil }
func (fakeStmt) Query([]driver.Value) (driver.Rows, error) {
	return &fakeRows{}, nil
}

func (fakeStmt) Exec([]driver.Value) (driver.Result, error) {
	return nil, driver.ErrSkip
}

type fakeRows struct {
	called bool
}

func (r *fakeRows) Columns() []string { return []string{"DEP", "AGE"} }

func (r *fakeRows) Close() error { return nil }

// Next returns one good row, then a non-EOF error. database/sql treats any
// non-EOF error from the driver as a fatal iteration error, which rows.Err()
// exposes after rows.Next() returns false.
func (r *fakeRows) Next(dest []driver.Value) error {
	if !r.called {
		r.called = true
		dest[0] = "Engineering"
		dest[1] = 30
		return nil
	}
	return errors.New("scan error: simulated failure reading row 2")
}

func main() {
	db, err := sql.Open("fake", "dsn")
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("failed to close db: %v\n", err)
		}
	}()

	ctx := context.Background()

	dep, age, err := get1(ctx, db, "1")
	fmt.Printf("get1 (mistake): dep=%q age=%d err=%v\n", dep, age, err) // get1 (mistake): dep="Engineering" age=30 err=<nil>

	dep, age, err = get2(ctx, db, "1")
	fmt.Printf("get2 (fixed):   dep=%q age=%d err=%v\n", dep, age, err) // get2 (fixed):   dep="" age=0 err=scan error: simulated failure reading row 2
}
