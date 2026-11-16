# httptest example

This example demonstrates Go's [`net/http/httptest`](https://pkg.go.dev/net/http/httptest) utility package for testing HTTP code without making real network calls.

## What's covered

**`main.go`**

* `Handler` — an HTTP handler that reads the request body, prefixes it with `hello `, and responds with a `201 Created` status plus an `X-API-VERSION` header.
* `DurationClient` — an HTTP client that POSTs a JSON request body to a URL and parses a `{"duration": N}` JSON response into a `time.Duration`.

**`main_test.go`**

* `TestHandler` — tests the handler in-process using `httptest.NewRequest` and `httptest.NewRecorder`. No real server is started; the request is dispatched directly to the handler.
* `TestDurationClientGet` — tests the client against a throwaway server created with `httptest.NewServer`. The client talks to `srv.URL` instead of a real external endpoint, so the test is fast, deterministic, and works offline.

## Running

```bash
go test -v
```

## Why use `httptest`?

* **Isolation** — tests don't depend on external services or network availability.
* **Speed** — no real TCP connections or DNS lookups.
* **Determinism** — you control the exact response the server returns.
* **No port conflicts** — `httptest.NewServer` picks a free port automatically.

See the [Go testing guide on httptest](https://go.dev/doc/tutorial/web-service-testing) for more on testing HTTP handlers and clients.
