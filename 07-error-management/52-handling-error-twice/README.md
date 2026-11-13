# Handling an Error Twice

Example from **Mistake 52: Handling an error twice** in *100 Go Mistakes and How to Avoid Them*.

## The mistake

An error should be handled **exactly once**. A common mistake is to both
log the error and return it:

```go
err := validateCoordinates(srcLat, srcLng)
if err != nil {
    log.Println("failed to validate source coordinates") // handles the error
    return Route{}, err                                   // handles it again
}
```

This produces duplicate log messages (once here, once at the caller) and
leaves the caller unsure whether it is responsible for logging. Conversely,
silently swallowing the error without logging or returning it loses
information entirely.

## The examples

The file `main.go` shows three ways to validate route coordinates:

### `GetRoute1` — ❌ Handle the error twice

`validateCoordinates1` both logs (`log.Printf`) and returns the error, and
`GetRoute1` logs again before returning. The same error is logged multiple
times at different levels of the call stack.

### `GetRoute2` — ⚠️ Returns without logging

`validateCoordinates2` only returns the error, and `GetRoute2` only returns it. Neither logs. `main()` logs the returned error at the boundary, so operators do get visibility — but `GetRoute2` itself adds no context, so the log message is just the raw validation error.

### `GetRoute3` — ✅ Handle the error once, at the right level

`validateCoordinates2` returns the error without logging. `GetRoute3`
wraps it with additional context using `%w` and returns it:

```go
err := validateCoordinates2(srcLat, srcLng)
if err != nil {
    return Route{},
        fmt.Errorf("failed to validate source coordinates: %w", err)
}
```

The error is handled exactly once: it is annotated with context and passed
up the call stack. The ultimate caller decides whether to log it, return it
to an HTTP handler, or retry.

## Rule of thumb

- **Log it** — if you are the final handler (e.g., in `main` or an HTTP handler).
- **Return it** — otherwise, optionally wrapping it with `%w` to add context.
- Don't do both for the same error.

## Running the example

```sh
go run main.go
```

`main.go` includes a `main()` that calls each `GetRoute*` function with valid and invalid coordinates. For invalid coordinates, `main()` logs once at the boundary — this is the correct place to log. Notice the difference:

- `GetRoute1` logs the same error **multiple times** (inside `validateCoordinates1`, inside `GetRoute1`, and again in `main()`) — the anti-pattern.
- `GetRoute2` logs **once** (in `main()`), but with no added context.
- `GetRoute3` logs **once** (in `main()`), with context wrapped in via `%w`.
