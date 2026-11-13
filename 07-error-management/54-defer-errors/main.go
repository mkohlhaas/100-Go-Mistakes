package main

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"log"
)

const query = "SELECT balance FROM clients WHERE id = ?"

func getBalance1(db *sql.DB, clientID string) (float32, error) {
	rows, err := db.Query(query, clientID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// Use rows
	return 0, nil
}

func getBalance2(db *sql.DB, clientID string) (float32, error) {
	rows, err := db.Query(query, clientID)
	if err != nil {
		return 0, err
	}
	defer func() { _ = rows.Close() }()

	// Use rows
	return 0, nil
}

func getBalance3(db *sql.DB, clientID string) (balance float32, err error) {
	rows, err := db.Query(query, clientID)
	if err != nil {
		return 0, err
	}
	defer func() {
		closeErr := rows.Close()
		if err != nil {
			if closeErr != nil {
				log.Printf("failed to close rows: %v", closeErr)
			}
			return
		}
		err = closeErr
	}()

	// Use rows
	return 0, nil
}

// ------------------------------------------------------------------ //
// ---- Dummy Driver so the Example Runs without a Real Database ---- //
// ------------------------------------------------------------------ //

type dummyDriver struct{}

func (dummyDriver) Open(name string) (driver.Conn, error) {
	return dummyConn{}, nil
}

type dummyConn struct{}

func (dummyConn) Prepare(query string) (driver.Stmt, error) {
	return dummyStmt{}, nil
}
func (dummyConn) Close() error              { return nil }
func (dummyConn) Begin() (driver.Tx, error) { return nil, errors.New("not supported") }

type dummyStmt struct{}

func (dummyStmt) Close() error  { return nil }
func (dummyStmt) NumInput() int { return 1 }
func (dummyStmt) Query(args []driver.Value) (driver.Rows, error) {
	return dummyRows{}, nil
}

func (dummyStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.New("not supported")
}

type dummyRows struct{}

func (dummyRows) Columns() []string { return []string{"balance"} }
func (dummyRows) Close() error      { return errors.New("close failed") }
func (dummyRows) Next(dest []driver.Value) error {
	return io.EOF
}

func init() {
	sql.Register("dummy", dummyDriver{})
}

// -------------- //
// ---- Main ---- //
// -------------- //

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := sql.Open("dummy", "")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer func() { _ = db.Close() }()

	clientID := "client-42"

	for _, fn := range []struct {
		name string
		f    func(*sql.DB, string) (float32, error)
	}{
		{"getBalance1", getBalance1}, // 2026/11/13 20:40:49 main.go:114: getBalance1: balance=0, err=<nil>
		{"getBalance2", getBalance2}, // 2026/11/13 20:40:49 main.go:114: getBalance2: balance=0, err=<nil>
		{"getBalance3", getBalance3}, // 2026/11/13 20:40:49 main.go:114: getBalance3: balance=0, err=close failed
	} {
		balance, err := fn.f(db, clientID)
		log.Printf("%s: balance=%v, err=%v", fn.name, balance, err)
	}
}
