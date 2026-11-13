# Comparing Error Types

This example demonstrates a common mistake in Go error handling: comparing errors with a **type switch**, which breaks as soon as errors are wrapped with `fmt.Errorf(..., "%w")`.

## The code

`main.go` contains two HTTP handlers that treat "transient" errors differently from other errors:

- `handler1` uses a **type switch** (`switch err := err.(type)`) to detect `transientError`.
- `handler2` uses **`errors.As`**, which walks the error chain and finds `transientError` even when it is wrapped.

Both handlers call `getTransactionAmount*` helpers that wrap the DB error with `fmt.Errorf(..., "%w")`. Because of that wrapping:

```go
switch err := err.(type) {
case transientError:
    http.Error(w, err.Error(), http.StatusServiceUnavailable)
default:
    http.Error(w, err.Error(), http.StatusBadRequest)
}
```

A type switch only matches the **top-level** error. If a lower layer wraps the error:

```go
return 0, fmt.Errorf("failed to get transaction %s: %w", transactionID, err)
```

then `err` is no longer a `transientError` — it's a `*fmt.wrapError`. The type switch falls through to `default`, and a transient failure is incorrectly reported as a `400 Bad Request` instead of a `503 Service Unavailable`.

## The fix

Use `errors.As`, which unwraps the error chain:

```go
if errors.As(err, &transientError{}) {
    http.Error(w, err.Error(), http.StatusServiceUnavailable)
} else {
    http.Error(w, err.Error(), http.StatusBadRequest)
}
```

`errors.As` returns `true` if `err` (or any error it wraps) matches the target type, so wrapped transient errors are still handled correctly.

### `errors.Is` vs. `errors.As`

Both walk the wrapped error chain, but they answer different questions:

| | `errors.Is(err, target)` | `errors.As(err, &target)` |
|---|---|---|
| **Question** | "Does this error equal this sentinel?" | "Does this error match this type?" |
| **Target** | An `error` **value** (usually a sentinel variable) | A pointer to a **type** or interface |
| **Returns** | `true` / `false` | `true` / `false`, plus the matching value |

Use `errors.Is` for sentinel errors (`io.EOF`, `os.ErrNotExist`, package-level `var ErrX`):

```go
if errors.Is(err, os.ErrNotExist) { ... }
```

Use `errors.As` for custom error types when you need to extract typed information:

```go
var perr *os.PathError
if errors.As(err, &perr) {
    fmt.Println("failed op:", perr.Op)
}
```

> Prefer `errors.Is` for sentinel errors and `errors.As` for custom error types. Avoid type switches for error comparison — they don't survive wrapping.

## Run it

```bash
go run main.go
```

## Related

- [Go blog: Go 1.13 error handling](https://go.dev/blog/go1.13-errors) — `errors.Is`, `errors.As`, and `%w`
- [100 Go Mistakes — Mistake 50: Compare error type](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them)