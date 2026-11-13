package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func listing1(err error) {
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("listing1: detected no rows via ==")
		} else {
			fmt.Println("listing1: missed no rows (err wrapped, == fails)")
		}
	}
}

func listing2(err error) {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("listing2: detected no rows via errors.Is")
		} else {
			fmt.Println("listing2: missed no rows")
		}
	}
}

func query(ok bool) error {
	if !ok {
		return nil
	}

	// Simulate a driver wrapping sql.ErrNoRows with %w.
	return fmt.Errorf("query failed: %w", sql.ErrNoRows)
}

func main() {
	listing1(query(true)) // listing1: missed no rows (err wrapped, == fails)
	listing2(query(true)) // listing2: detected no rows via errors.Is
}
