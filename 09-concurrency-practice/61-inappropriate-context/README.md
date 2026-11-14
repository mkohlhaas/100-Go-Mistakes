# 61 - Inappropriate Context

This directory contains a Go snippet demonstrating **Mistake 61: using an inappropriate `context.Context` when spawning goroutines that outlive an HTTP request** (from *100 Go Mistakes and How to Avoid Them*).

## Problem

When an HTTP handler spawns a background goroutine, the context passed to that goroutine decides:

- whether the goroutine keeps running after the handler returns,
- whether request-scoped values (tracing IDs, user info, …) are still available,
- whether the server cancels the work when the client disconnects.

## Handlers

- `handler1` — spawns the goroutine with `r.Context()` (the request context). The server cancels this context when the request finishes, so work started in the goroutine can be cancelled mid-flight.
- `handler2` — spawns the goroutine with `context.Background()`. The work outlives the request, but every request-scoped value, deadline, and cancellation signal is lost.
- `handler3` — spawns the goroutine with a `detach` wrapper: a custom `context.Context` that never gets cancelled (its `Done()` channel is `nil`) but still forwards `Value` lookups to the original request context. The work outlives the request while keeping request-scoped values.

## Structure

- `main.go` — main package with the three handlers plus a `main()` that starts an HTTP server on `:8080` and shuts down gracefully on `SIGINT`/`SIGTERM`.

## How to run

```bash
go run .
```

Then hit the endpoints:

```bash
curl http://localhost:8080/handler1
curl http://localhost:8080/handler2
curl http://localhost:8080/handler3
```

Stop the server with `Ctrl+C`.

## What this code demonstrates

| Handler | Context passed to goroutine | Outlives request? | Keeps request-scoped values? | Cancellable with request? |
|---|---|---|---|---|
| `handler1` | `r.Context()` | ❌ (cancelled when the request ends) | ✅ | ✅ |
| `handler2` | `context.Background()` | ✅ | ❌ | ❌ |
| `handler3` | `detach` wrapper | ✅ | ✅ | ❌ |

## Guidance

- Don't pass `r.Context()` to background goroutines that must outlive the HTTP request — the server cancels it when the request finishes.
- Don't blindly use `context.Background()` in request handlers — you lose tracing, deadlines, and request-scoped values.
- If background work needs to outlive the request but still carry request-scoped values, detach cancellation from the context while forwarding `Value` (as `handler3` does), or use a dedicated service context.

## Related

Part of the [100 Go Mistakes](https://github.com/teivah/100-go-mistakes) series — mistake 61: inappropriate context.
