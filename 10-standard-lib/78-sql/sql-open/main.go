package main

import (
	"database/sql"
	"fmt"
	"log"

	// Register the MySQL driver.
	_ "github.com/go-sql-driver/mysql"
)

var dsn = ""

func listing1() error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error while opening the database.")
		return err
	}

	if err := db.Ping(); err != nil {
		fmt.Println("Error while pinging the database.") // Error while pinging the database.
		return err
	}

	_ = db
	return nil
}

func main() {
	if err := listing1(); err != nil {
		log.Printf("listing1 failed: %v", err) // 2026/11/15 16:11:49 listing1 failed: dial tcp 127.0.0.1:3306: connect: connection refused
		return
	}
	fmt.Println("listing1 succeeded")
}
