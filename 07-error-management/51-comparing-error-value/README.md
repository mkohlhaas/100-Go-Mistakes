# 51 — Comparing error values

Example code for mistake #51 from *100 Go Mistakes*: comparing errors with `==`
instead of using `errors.Is`.

## The mistake

`listing1` compares an error value directly with `==`:

```go
if err == sql.ErrNoRows { ... }
```

This only works when `err` is *exactly* `sql.ErrNoRows`. If a database driver
wraps the error (e.g. `fmt.Errorf("query failed: %w", sql.ErrNoRows)`), the
comparison fails and the caller never detects the "no rows" condition.

## The fix

`listing2` uses `errors.Is`, which walks the error chain and unwraps wrapped
errors:

```go
if errors.Is(err, sql.ErrNoRows) { ... }
```

Prefer `errors.Is` whenever you need to check for a sentinel error.

## Run it

```bash
go run main.go
```

Both `listing1` and `listing2` compile and run, but only `listing2` correctly
detects sentinel errors when they are wrapped.

## See also

- [`errors.Is`](https://pkg.go.dev/errors#Is)
- [`fmt.Errorf` with `%w`](https://pkg.go.dev/fmt#hdr-Error_formatting)