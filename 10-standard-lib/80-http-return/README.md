# Mistake #80: Forgetting to `return` after `http.Error` in a handler

When an HTTP handler encounters an error it typically calls [`http.Error`](https://pkg.go.dev/net/http#Error) to send an error status and body to the client. The mistake is forgetting to `return` immediately after that call. Execution then falls through, and the handler writes more data to the `ResponseWriter` — sometimes even calling `WriteHeader` a second time.

This produces an invalid HTTP response: headers are already on their way to the client, so any subsequent write corrupts the response, and Go's HTTP server logs a `superfluous response.WriteHeader call` warning.

## The bug

```go
func handler1(w http.ResponseWriter, req *http.Request) {
	err := foo(req)
	if err != nil {
		http.Error(w, "foo", http.StatusInternalServerError)
	}

	_, _ = w.Write([]byte("all good"))
	w.WriteHeader(http.StatusCreated)
}
```

If `foo` returns an error:

1. `http.Error` sends a `500 Internal Server Error` (headers + body).
2. Execution continues — `w.Write([]byte("all good"))` appends more body bytes after the error body.
3. `w.WriteHeader(http.StatusCreated)` tries to change the status to `201` after headers were already sent: a no-op that Go flags as superfluous.

The client receives a malformed response mixing the error and the success payload.

## The fix

Always `return` right after writing an error response:

```go
func handler2(w http.ResponseWriter, req *http.Request) {
	err := foo(req)
	if err != nil {
		http.Error(w, "foo", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("all good"))
	w.WriteHeader(http.StatusCreated)
}
```

Now, on error, the handler stops immediately after `http.Error` and never touches the response again.

## Why it matters

- `http.Error` commits the response (status line, headers, body) to the client.
- Anything written to the `ResponseWriter` after that point is part of the same HTTP response and corrupts it.
- A `WriteHeader` call after headers have been sent is ignored and logged as a warning.

The rule for `net/http` handlers: once you start writing the response, finish it — and on an error path, that means `return` after `http.Error`.

## Running the code

The handlers are exercised through unit tests (see below). There is no `main()` or HTTP server to run — `go test` is the entry point:

```bash
go test -v
```

This compiles the package and runs `TestHandlerBuggy` and `TestHandlerFixed`, which invoke `handler1` and `handler2` directly through [`net/http/httptest`](https://pkg.go.dev/net/http/httptest) and assert on the resulting status codes and bodies.

## Testing with `httptest`

You don't need a running server to exercise the handlers. Go's [`net/http/httptest`](https://pkg.go.dev/net/http/httptest) package lets you capture responses in a unit test:

```go
package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerBuggy(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/buggy?err=1", nil)
	rec := httptest.NewRecorder()

	handler1(rec, req)

	res := rec.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("failed to close response body: %v", err)
		}
	}()
	body, _ := io.ReadAll(res.Body)

	// The bug: the response body contains both the error and the success payload.
	if got := string(body); !strings.Contains(got, "foo") || !strings.Contains(got, "all good") {
		t.Fatalf("expected corrupted body mixing error and success, got %q", got)
	}
	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.StatusCode)
	}
}

func TestHandlerFixed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/fixed?err=1", nil)
	rec := httptest.NewRecorder()

	handler2(rec, req)

	res := rec.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("failed to close response body: %v", err)
		}
	}()
	body, _ := io.ReadAll(res.Body)

	// The fix: only the error body is present, and the status is 500.
	if got := string(body); strings.Contains(got, "all good") {
		t.Fatalf("expected no success payload after error, got %q", got)
	}
	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.StatusCode)
	}
}
```

Run it with:

```bash
go test -v
```

This is the idiomatic way to test `net/http` handlers: build a request with `httptest.NewRequest`, capture the response with `httptest.NewRecorder`, invoke the handler directly, then assert on the status code and body.

---

Part of the [100 Go Mistakes](https://100go.co) repository.