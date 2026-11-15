# prepared-statements

Example for mistake **78** from *100 Go Mistakes*: leaking database resources
by never closing prepared statements (and the rows they return).

## The mistake

`listing1` prepares a statement and queries it, but never calls
`stmt.Close()` or `rows.Close()`. Every call leaks a server-side prepared
statement and a result set until the connection is discarded or closed.

```go
stmt, err := db.Prepare("SELECT * FROM orders WHERE id = ?")
if err != nil {
    return err
}
// BUG: stmt is never closed
rows, err := stmt.Query(id)
if err != nil {
    return err
}
// BUG: rows is never closed
_ = rows
return nil
```

## The fix

Close both resources as soon as they are obtained, typically with `defer`
right after the error check:

```go
stmt, err := db.Prepare("SELECT * FROM orders WHERE id = ?")
if err != nil {
    return err
}
defer stmt.Close()

rows, err := stmt.Query(id)
if err != nil {
    return err
}
defer rows.Close()

_ = rows
return nil
```

See `listing2` in `main.go` for the corrected version.

### Do we have to close everything?

Yes — close all three, but they serve different purposes:

| Resource | Why close it |
|----------|--------------|
| `*sql.Rows` | **Required.** Holds a connection from the pool until closed. Not closing it leaks connections and can exhaust the pool. |
| `*sql.Stmt` | **Recommended.** Frees server-side prepared statement resources promptly. |
| `*sql.DB` | **Required when done.** Closes the entire connection pool (e.g., on app shutdown). |

`rows.Close()` is the most critical — leaving rows open ties up a database connection
indefinitely.

`Close()` can also return an error (e.g., when there are still open statements or
rows on the connection, or due to a network problem), so the example wraps each
`defer` in an anonymous function and checks the error:

```go
defer func() {
    if err := db.Close(); err != nil {
        log.Printf("closing db: %v", err)
    }
}()
```

## Running

```bash
go run .
```

The program opens an in-memory SQLite database, creates an `orders` table, and
runs both the leaky (`listing1`) and the corrected (`listing2`) versions.

## Notes

* The book snippet uses the table name `ORDER`, which is a reserved keyword in
  SQLite, so this runnable example uses `orders` instead.
* Driver: [`github.com/mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3)
  (CGO, requires a C compiler such as `gcc`).
* `go.mod` targets Go 1.24.
