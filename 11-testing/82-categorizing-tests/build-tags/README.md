# Categorizing Tests with Build Tags

This example shows how to separate **integration tests** from **unit tests** using Go build tags, so integration tests (which may require a database, network, or other external services) only run when explicitly requested.

## Files

```
db.go             # library under test
db_unit_test.go   # unit tests — no build tag, always run
db_test.go        # integration tests — gated behind the `integration` build tag
go.mod            # module definition
```

## The example

`db_test.go` contains tests gated behind the `integration` build tag:

```go
//go:build integration
// +build integration

package db

func TestInsert1(t *testing.T) {
    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        t.Skip("DB_DSN not set; skipping integration test")
    }
    // ... insert against a real database
}

func TestInsert2(t *testing.T) {
    if os.Getenv("INTEGRATION") != "true" {
        t.Skip("skipping integration test")
    }
    // ... insert against a real database
}
```

Two mechanisms are combined:

1. **Build tag `integration`** — The file is only compiled when you pass `-tags integration` to `go test`. Without the tag, these tests are completely excluded from the build and won't run (or even compile).
2. **Environment guard `INTEGRATION=true`** — Even when the file is compiled, `TestInsert2` skips unless `INTEGRATION=true` is set, giving a second safety net so tests never accidentally run against a shared/staging database.

## Running the tests (you can try this now)

### Unit tests only (default)

```bash
go test ./...
```

Output:

```
ok  	example.com/categorizing-tests	0.003s
```

Only the four unit tests in `db_unit_test.go` run. `db_test.go` is not compiled at all.

### Integration tests compiled but skipped

```bash
go test -tags integration ./... -v
```

Output:

```
=== RUN   TestInsert1
    db_test.go:16: DB_DSN not set; skipping integration test
--- SKIP: TestInsert1 (0.00s)
=== RUN   TestInsert2
    db_test.go:34: skipping integration test
--- SKIP: TestInsert2 (0.00s)
=== RUN   TestOpenEmptyDSN
--- PASS: TestOpenEmptyDSN (0.00s)
...
```

Now `db_test.go` **is** compiled, but both integration tests skip because neither `DB_DSN` nor `INTEGRATION=true` is set. You can see they were considered and deliberately skipped rather than silently ignored.

### Actually running integration tests

Set the required environment variables (point `DB_DSN` at a real database you own):

```bash
INTEGRATION=true DB_DSN="host=localhost port=5432 dbname=mydb" \
  go test -tags integration ./...
```

### CI recommendation

Run unit tests on every commit/PR:

```bash
go test ./...
```

Run integration tests separately (e.g., nightly or in a dedicated CI job) with the build tag and environment variables set:

```bash
INTEGRATION=true DB_DSN=... go test -tags integration ./...
```

## Why this matters

- Keeps the fast unit-test suite fast — no accidental database or network calls.
- Prevents integration tests from running against shared environments by mistake.
- Makes it obvious which tests are integration tests by build tag rather than by convention alone.

## Go version note

The file uses both the modern `//go:build integration` syntax (Go 1.17+) and the legacy `// +build integration` syntax for compatibility with older Go toolchains. Keep the two lines in sync when editing.