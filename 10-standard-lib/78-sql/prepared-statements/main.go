package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// listing1 demonstrates the mistake: the prepared statement (and the rows it
// returns) are never closed, leaking database resources.
//
// Note: the original snippet uses the table name ORDER, which is a reserved
// keyword in SQLite, so this runnable example uses orders instead.
func listing1(db *sql.DB, id string) error {
	stmt, err := db.Prepare("SELECT * FROM orders WHERE id = ?")
	if err != nil {
		return err
	}
	// BUG: stmt is never closed.
	rows, err := stmt.Query(id)
	if err != nil {
		return err
	}
	// BUG: rows is never closed.
	_ = rows
	return nil
}

// listing2 shows the correct way: always close the statement and the rows,
// ideally with defer right after checking for an error.
func listing2(db *sql.DB, id string) error {
	stmt, err := db.Prepare("SELECT * FROM orders WHERE id = ?")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("closing stmt: %v", err)
		}
	}()

	rows, err := stmt.Query(id)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("closing rows: %v", err)
		}
	}()

	_ = rows
	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("closing db: %v", err)
		}
	}()

	if _, err := db.Exec("CREATE TABLE orders (id INTEGER PRIMARY KEY, name TEXT)"); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec("INSERT INTO orders (name) VALUES ('apple'), ('banana')"); err != nil {
		log.Fatal(err)
	}

	if err := listing1(db, "1"); err != nil {
		log.Printf("listing1: %v", err)
	} else {
		fmt.Println("listing1: ran (statement/rows leaked!)")
	}

	if err := listing2(db, "1"); err != nil {
		log.Printf("listing2: %v", err)
	} else {
		fmt.Println("listing2: ran (resources closed)")
	}
}
