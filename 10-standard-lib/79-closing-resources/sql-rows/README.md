# sql-rows — Closing `sql.Rows` (Go Mistake #79)

This example demonstrates a common mistake when querying a database with Go's
`database/sql` package: **forgetting to close `sql.Rows`**.

An open `sql.Rows` holds a connection from the database's connection pool.
Until `rows.Close()` is called, that connection stays checked out and cannot be
reused. Forgetting to close rows leaks connections and eventually exhausts the
pool, so every query must ensure `rows.Close()` is called — even on early
returns and error paths.

## The code

```go
package main

import (
	"database/sql"
	"log"
)

func listing1() error {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT * FROM CUSTOMERS")
	if err != nil {
		return err
	}

	// Use rows
	_ = rows

	return nil
}

func listing2() error {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT * FROM CUSTOMERS")
	if err != nil {
		return err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v\n", err)
		}
	}()

	// Use rows
	_ = rows

	return nil
}

var dataSourceName = ""
```

## What happens

### `listing1` — rows are never closed (resource leak)

`db.Query` returns an open `sql.Rows` backed by a pooled connection. The
function returns without calling `rows.Close()`, so the connection remains
checked out from the pool indefinitely. Repeated calls exhaust the connection
pool and subsequent queries block or fail.

### `listing2` — rows are closed with `defer` (correct)

`rows.Close()` is registered with `defer` immediately after the error from
`db.Query` is checked. When the function returns — normally or early — the
deferred call releases the connection back to the pool. A close failure is
logged rather than ignored.

## The lesson

When using `database/sql`:

1. **Always call `rows.Close()`** — register it with `defer` right after
   checking the error from `db.Query`/`db.QueryRowContext`, before doing
   anything else with the rows. The connection stays reserved until close.
2. **Close rows even on early returns** — `defer` guarantees the close runs no
   matter how the function exits.
3. **Check `rows.Err()` after iterating** — after a `for rows.Next()` loop,
   call `rows.Err()` to surface any error that occurred during iteration
   (e.g. a connection drop mid-stream). `rows.Close()` does not report it.
4. **`rows.Close()` is safe to call multiple times** — calling it on already
   closed rows is a no-op.

The correct pattern:

```go
rows, err := db.Query("SELECT * FROM CUSTOMERS")
if err != nil {
    return err
}
defer rows.Close()

for rows.Next() {
    // scan each row
}
if err := rows.Err(); err != nil {
    return err
}
```

Note that `defer` only defers the *call*; the error from `rows.Close()` should
still be handled — either logged as above, or combined with other errors using
`errors.Join` — so a failed close is not silently swallowed.

## Run it

This directory is part of the root module (`module 100-go-mistakes`), so there's
no local `go.mod`. Build and run from this directory:

```bash
go run .
```

The example compiles but does not connect to a real database
(`dataSourceName` is empty), so running it won't perform a query — it exists to
illustrate the pattern.

## Reference

*Mistake #79: Closing resources* from *100 Go Mistakes and How to Avoid Them*.

- Sibling examples: `../http/` (HTTP response bodies), `../os-file/` (file handles)
- Related: `../../78-sql/prepared-statements/` covers the same rows-closing idea in the context of prepared statements (Mistake #78), including `stmt.Close()` and a runnable `main()` with an in-memory SQLite database.
