# 60 - Contexts

This directory contains a Go snippet demonstrating the use of `context.Context` for cancellation, timeouts, and request-scoped values.

## Structure

- `main.go` — main package with helpers showing context usage:
  - `publishPosition` — creates a context with a 4-second timeout
  - `checkValid` — HTTP middleware storing a request-scoped value via `context.WithValue`
  - `handler` — select loop that exits on `ctx.Done()`
- `flight/flight.go` — `flight.Position` struct used by the publisher interface
- `go.mod` — module definition (`module contexts`)

## How to run

This directory is its own Go module and includes a `main()` function, so it can be built and run directly:

```bash
go mod tidy
go run .
```

Sample output:

```
publishPosition succeeded
slowPublisher timed out: context deadline exceeded
manual cancel returned: context canceled
handler returned: context canceled
```

Or build and vet it:

```bash
go build ./...
go vet ./...
```

## Import path

The `flight` package is imported from the local subdirectory:

```go
import "contexts/flight"
```

## What this code demonstrates

This snippet demonstrates the **correct way to use `context.Context` in Go** (Mistake 60 from *100 Go Mistakes*). It covers three core patterns:

| Pattern | Code | Purpose |
|---|---|---|
| **Timeouts & cancellation** | `publishPosition` with `context.WithTimeout(..., 4*time.Second)` + `defer cancel()` | Sets a deadline on an operation and releases resources when done, preventing context leaks |
| **Request-scoped values** | `checkValid` with `context.WithValue(r.Context(), isValidHostKey, validHost)` | Carries per-request data through middleware/handlers |
| **Cancellation propagation** | `handler` with `<-ctx.Done()` / `return ctx.Err()` | A long-running loop respects cancellation and exits cleanly when the context is done |
| **Checking ctx.Done()** | `mockPublisher.Publish` / `slowPublisher.Publish` selecting on `<-ctx.Done()` vs `<-time.After(...)` | A callee actively checks whether the caller cancelled the context and aborts early; `slowPublisher` (5s work) against a 4s timeout returns `context.DeadlineExceeded` |

Supporting details:

- `type key string` / `isValidHostKey` — a private key type avoids collisions with string keys from other packages.
- Context is passed as a function parameter through `Publish(ctx, ...)` and `publishHandler`, not stored in a struct — context travels through the call stack.

## Is `case <-ctx.Done():` typical in functions with `context.Context`?

**Yes — but only in functions that do blocking or long-running work.**

| If the function... | Check `ctx.Done()`? | Example |
|---|---|---|
| Does blocking work (loops, I/O, sleep, waiting on channels) | **Yes** | `handler`, `slowPublisher.Publish` |
| Just passes `ctx` to another function | No — let the callee check it | `publishPosition` passing `ctx` to `Publish` |
| Returns immediately | No | simple getter/setter |

### Typical pattern

```go
func doWork(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err() // propagate cancellation/timeout
        default:
            // do a chunk of work
        }
    }
}
```

Or in a `select` with other channels:

```go
select {
case msg := <-ch:
    // process msg
case <-ctx.Done():
    return ctx.Err()
}
```

### Why this matters

If you don't check `ctx.Done()`, your function ignores the caller's request to stop — for example, when an HTTP request is cancelled by the client, or a timeout fires. That causes resource leaks (goroutines stuck waiting, DB connections held open, etc.).

Rule of thumb: **if your function blocks for any meaningful time, include `case <-ctx.Done():` and return `ctx.Err()`**.

## Related

Part of the [100 Go Mistakes](https://github.com/teivah/100-go-mistakes) series — mistake 60: contexts.
