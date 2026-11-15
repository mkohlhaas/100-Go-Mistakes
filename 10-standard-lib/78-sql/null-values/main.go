package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log"
)

func listing1(db *sql.DB, id string) error {
	rows, err := db.Query("SELECT DEP, AGE FROM EMP WHERE ID = ?", id)
	if err != nil {
		return err
	}
	// Defer closing rows

	var (
		department string // ⚠️
		age        int    // ⚠️
	)
	for rows.Next() {
		err := rows.Scan(&department, &age)
		if err != nil {
			return err
		}
		// ...
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func listing2(db *sql.DB, id string) error {
	rows, err := db.Query("SELECT DEP, AGE FROM EMP WHERE ID = ?", id)
	if err != nil {
		return err
	}
	// Defer closing rows

	var (
		department *string // ✅
		age        int     // ⚠️
	)
	for rows.Next() {
		err := rows.Scan(&department, &age)
		if err != nil {
			return err
		}
		// ...
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func listing3(db *sql.DB, id string) error {
	rows, err := db.Query("SELECT DEP, AGE FROM EMP WHERE ID = ?", id)
	if err != nil {
		return err
	}
	// Defer closing rows

	var (
		department sql.NullString // ✅
		age        sql.NullInt64  // ✅
	)
	for rows.Next() {
		err := rows.Scan(&department, &age)
		if err != nil {
			return err
		}
		// ...
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

// ---------------------------------------------------------------------------
// Dummy in-memory database/sql driver (standard library only).
//
// It lets us run listing1/listing2/listing3 against a table with SQL NULLs
// without installing a real database driver.
//
// Table EMP (ID, DEP, AGE):
//   ID=1: DEP="Engineering", AGE=NULL   -> listing1/listing2 error on plain int AGE
//   ID=2: DEP=NULL,      AGE=30        -> listing1 errors on plain string DEP,
//                                         listing2 handles nil *string, listing3 ok
//   ID=3: DEP="Sales",   AGE=25        -> all three ok
// ---------------------------------------------------------------------------

var empRows = []struct {
	id  string
	dep any // string or nil
	age any // int64 or nil
}{
	{"1", "Engineering", nil},
	{"2", nil, int64(30)},
	{"3", "Sales", int64(25)},
}

type dummyDriver struct{}

func (dummyDriver) Open(name string) (driver.Conn, error) { return dummyConn{}, nil }

type dummyConn struct{}

func (dummyConn) Prepare(query string) (driver.Stmt, error) {
	return dummyStmt{query: query}, nil
}
func (dummyConn) Close() error { return nil }
func (dummyConn) Begin() (driver.Tx, error) {
	return dummyTx{}, nil
}

type dummyTx struct{}

func (dummyTx) Commit() error   { return nil }
func (dummyTx) Rollback() error { return nil }

type dummyStmt struct{ query string }

func (dummyStmt) Close() error  { return nil }
func (dummyStmt) NumInput() int { return 1 }

func (s dummyStmt) Query(args []driver.Value) (driver.Rows, error) {
	id := ""
	if len(args) > 0 {
		id = fmt.Sprint(args[0])
	}
	var rows []struct {
		dep any
		age any
	}
	for _, e := range empRows {
		if e.id == id {
			rows = append(rows, struct {
				dep any
				age any
			}{dep: e.dep, age: e.age})
		}
	}
	return &dummyRows{rows: rows}, nil
}

func (dummyStmt) Exec(args []driver.Value) (driver.Result, error) {
	return dummyResult{}, nil
}

type dummyResult struct{}

func (dummyResult) LastInsertId() (int64, error) { return 0, nil }
func (dummyResult) RowsAffected() (int64, error) { return 0, nil }

type dummyRows struct {
	rows []struct {
		dep any
		age any
	}
	idx int
}

func (r *dummyRows) Columns() []string { return []string{"DEP", "AGE"} }

func (r *dummyRows) Close() error { return nil }

func (r *dummyRows) Next(dest []driver.Value) error {
	if r.idx >= len(r.rows) {
		return io.EOF
	}
	row := r.rows[r.idx]
	r.idx++
	dest[0] = row.dep
	dest[1] = row.age
	return nil
}

func main() {
	// Register the dummy driver so sql.Open can use it.
	sql.Register("dummy", dummyDriver{})

	db, err := sql.Open("dummy", "dummy://db")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("closing db: %v", cerr)
		}
	}()

	for _, id := range []string{"1", "2", "3"} {
		fmt.Printf("=== ID %s ===\n", id)

		fmt.Print("listing1 (plain string/int): ")
		if err := listing1(db, id); err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("ok")
		}

		fmt.Print("listing2 (*string/int): ")
		if err := listing2(db, id); err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("ok")
		}

		fmt.Print("listing3 (sql.Null*): ")
		if err := listing3(db, id); err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("ok")
		}
	}
}
