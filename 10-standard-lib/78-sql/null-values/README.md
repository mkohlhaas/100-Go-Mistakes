# Mistake 78 — SQL NULL values

Chapter 10: Standard library · [`100 Go Mistakes and How to Avoid Them`](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them) · [100go.co](https://100go.co)

Go database types (`string`, `int`, `int64`, `float64`, `bool`, `time.Time`) cannot hold SQL `NULL`. Scanning a `NULL` column into one of them makes `rows.Scan` return an error like `converting NULL to string is unsupported`.

## The mistake

`listing1` scans `DEP` into a plain `string` and `AGE` into a plain `int`:

```go
var (
    department string
    age        int
)
for rows.Next() {
    err := rows.Scan(&department, &age)
    if err != nil {
        return err
    }
}
```

If `DEP` or `AGE` is `NULL` in the database, `Scan` fails and the whole query aborts.

## The fix

Use the `database/sql` null-safe types, or pointers, so a `NULL` column doesn't break the scan:

### `sql.NullString` / `sql.NullInt64` / … (recommended)

`listing3` uses `sql.NullString` for `DEP`:

```go
var (
    department sql.NullString
    age        sql.NullInt64
)
for rows.Next() {
    err := rows.Scan(&department, &age)
    if err != nil {
        return err
    }
    if department.Valid {
        // use department.String
    }
    if age.Valid {
        // use age.Int64
    }
}
```

`database/sql` provides `NullString`, `NullInt64`, `NullInt32`, `NullInt16`, `NullFloat64`, `NullBool`, and `NullTime`. Each is a struct with a `Valid` bool plus a value field (`String`, `Int64`, `Float64`, `Bool`, `Time`). Always check `Valid` before reading the value — an invalid (NULL) value holds the zero value of the underlying type.

### Pointers

`listing2` uses `*string` for `DEP`; a `NULL` column becomes a `nil` pointer:

```go
var department *string
err := rows.Scan(&department, &age)
// department == nil when DEP is NULL
```

Pointers work but require nil checks on every column, and they don't exist for every SQL type. Prefer `sql.Null*` types when you want an explicit "present or absent" semantics.

## Running

`main.go` registers a dummy in-memory `database/sql` driver (standard library only,
no external packages or CGo) and runs `listing1`, `listing2`, and `listing3`
against a small `EMP` table with SQL `NULL`s:

| ID | DEP     | AGE   | listing1 (plain) | listing2 (*string) | listing3 (sql.Null*) |
|----|---------|-------|------------------|--------------------|----------------------|
| 1  | "Engineering" | NULL | ❌ error | ❌ error | ✅ ok |
| 2  | NULL    | 30    | ❌ error | ✅ ok | ✅ ok |
| 3  | "Sales" | 25    | ✅ ok | ✅ ok | ✅ ok |

```bash
go run .
```

Example output:

```
=== ID 1 ===
listing1 (plain string/int): error: sql: Scan error on column index 1, name "AGE": converting NULL to int is unsupported
listing2 (*string/int): error: sql: Scan error on column index 1, name "AGE": converting NULL to int is unsupported
listing3 (sql.Null*): ok
=== ID 2 ===
listing1 (plain string/int): error: sql: Scan error on column index 0, name "DEP": converting NULL to string is unsupported
listing2 (*string/int): ok
listing3 (sql.Null*): ok
=== ID 3 ===
listing1 (plain string/int): ok
listing2 (*string/int): ok
listing3 (sql.Null*): ok
```

To run against a real database instead, replace the dummy driver with a real one
(e.g. `_ "github.com/go-sql-driver/mysql"`), open a connection with
`sql.Open("<driver>", "<data-source-name>")`, and call `listing3(db, "1")`. The Go
module is resolved from the repository root (`go.mod`), so no `go.mod` is needed in
this directory.

## Key takeaway

Never scan a possibly-`NULL` column into a plain Go type. Use `sql.NullString`, `sql.NullInt64`, etc. (check `.Valid`), or a pointer type, so `NULL` rows are handled instead of breaking `rows.Scan`.
