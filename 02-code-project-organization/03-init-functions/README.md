# Mistake 3 — Init functions

This example demonstrates the pitfalls of using Go's `init()` functions for
package-level setup.

## The problem

`init()` runs automatically when a package is imported, before `main()` starts.
That makes it tempting to put connection setup, global state initialization, and
other side effects there — but it introduces several problems:

1. **Hidden side effects at import time.** Importing a package runs its
   `init()` even if the caller doesn't need the setup. The `db` package forces a
   MySQL connection (and panics on failure) just by being imported.

2. **`init()` cannot report errors.** It has no return values, so the only way
   to signal failure is to panic or log.Fatal. The `db` package's `init()`
   panics when `MYSQL_DATA_SOURCE_NAME` is missing or the database can't be
   reached.

3. **Implicit, fragile ordering.** `init()` ordering is determined by import
   dependencies and source-file order, which is hard to trace and easy to break
   when imports change.

4. **Hard to test.** Globals populated by `init()` can't be controlled per test
   or swapped for a mock.

## The better way

Replace `init()` with an explicit constructor function that returns errors:

```go
func createClient(dataSourceName string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dataSourceName)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
```

Call it from `main()`, handle the error, and pass the resulting client around
(dependency injection). This makes dependencies visible, failures handleable,
and the code testable.

## Files

- `db/main.go` — the mistake: an `init()` that opens a MySQL connection and
  panics on failure, alongside `createClient()`, the explicit alternative.
- `main/main.go` — demonstrates implicit `init()` ordering across packages.
- `redis/redis.go` — an `init()` that runs on import even though `main` only
  needs `Store`.
