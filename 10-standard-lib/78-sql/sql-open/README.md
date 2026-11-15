# 78 - `sql.Open` doesn't open a connection

This example illustrates a common mistake when working with Go's `database/sql` package: assuming that `sql.Open` actually opens a connection to the database and validates the data source name (DSN).

## The code

```go
package main

import "database/sql"

var dsn = ""

func listing1() error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error while opening the database.")
		return err
	}

	if err := db.Ping(); err != nil {
		fmt.Println("Error while pinging the database.")
		return err
	}

	_ = db
	return nil
}
```

## What happens

1. `sql.Open("mysql", dsn)` is called with an empty DSN.
2. It returns **no error** — `sql.Open` only checks that the driver name is registered; it does **not** connect to the database and does **not** validate the DSN.
3. The actual connection attempt happens lazily, the first time the database is used.
4. `db.Ping()` forces a real connection attempt and surfaces the error (invalid DSN, unreachable host, bad credentials, etc.).

## The lesson

`sql.Open` does not open a connection. It returns a `*sql.DB` handle (a connection pool) and only validates that the driver name is registered. The DSN is never checked until you perform an operation that requires a connection, such as `db.Ping()`, `db.Query()`, or `db.Exec()`.

To catch connection problems early, always follow `sql.Open` with `db.Ping()` and handle its error:

```go
db, err := sql.Open("mysql", dsn)
if err != nil {
    return err
}
defer db.Close()

if err := db.Ping(); err != nil {
    return fmt.Errorf("connect to database: %w", err)
}
```

Also keep in mind:

- `*sql.DB` is a pool safe for concurrent use. Create it once and reuse it; don't open and close it per request.
- `defer db.Close()` when the handle is no longer needed (short-lived programs, tests, etc.).
- Register the driver with a blank import: `_ "github.com/go-sql-driver/mysql"`.

## Run it

This directory is part of the root module (`module 100-go-mistakes`), so there's no local `go.mod`. A MySQL driver must be registered to build and run:

```bash
go run .
```

Sample output (no MySQL server running):

```
Error while pinging the database.
2026/11/15 16:15:13 listing1 failed: dial tcp 127.0.0.1:3306: connect: connection refused
```

Notice that `sql.Open` itself succeeded — the error only surfaced at `db.Ping()`, which is the first operation that actually opens a connection.

## Related

Part of the [100 Go Mistakes](https://github.com/teivah/100-go-mistakes) series — mistake 78: `sql.Open` doesn't open a connection.
