# Setup & Teardown

This package demonstrates Go test setup and teardown patterns against a **real** SQL database (SQLite via [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite)) — no Docker or external services required.

## Files

- `go.mod` — module definition with the SQLite driver dependency.
- `main_test.go` — working example of both patterns.

## Patterns demonstrated

| Pattern | Where | Use when |
|---|---|---|
| `TestMain(m *testing.M)` | `TestMain` | Expensive setup shared by every test in the package (run once). |
| `t.Cleanup(func ...)` | `createConnection` / `TestInsertThenCleanup` | Per-test cleanup that must run even if the test fails (e.g., close a connection, delete a temp row). |
| `defer` in `TestMain` | `defer db.Close()` | Cleanup tied to `TestMain`'s lifetime (note: runs only if `os.Exit` is not called first). |

## `TestMain`

`TestMain(m *testing.M)` is part of Go's standard library [`testing`](https://pkg.go.dev/testing) package.

If a test package contains a function with that exact signature, `go test` calls `TestMain` instead of running the tests directly. Inside `TestMain`, `m.Run()` actually executes all the tests in the package and returns an exit code. This lets you run setup before the tests and teardown after them:

```go
func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    teardown()
    os.Exit(code)
}
```

Note that `TestMain` must call `os.Exit(code)` itself — otherwise `go test` will ignore the result returned by `m.Run()`.

## Key details

- **Package-wide setup/teardown** (`setupDB`/`teardownDB`) runs once in `TestMain`: the schema is created and seed data inserted before `m.Run()`, and the database file is removed after all tests complete.
- **Per-test connection** (`createConnection`) opens its own `*sql.DB` and registers `conn.Close()` via `t.Cleanup`, so each test gets a fresh connection that is always closed.
- **Per-test data cleanup** (`TestInsertThenCleanup`) inserts a temporary row and removes it in a `t.Cleanup`, leaving the database clean for the next test.

## Running

```bash
go test -v ./...
```

## Adapting to a real database (MySQL, Postgres, …)

Swap the SQLite driver for your database driver and point `dsn` at your test instance. If your database runs in Docker, use [testcontainers-go](https://go.testcontainers.org/) inside `TestMain` to start/stop the container, and `t.Skip` when Docker is unavailable so the package still works in CI without a container runtime.
