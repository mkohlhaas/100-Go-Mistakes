# 79 - Closing HTTP response bodies

This example illustrates common mistakes when handling `http.Response` bodies with Go's `net/http` package: forgetting to close the body, closing it at the wrong time, and failing to drain it so the underlying connection can be reused.

## The code

```go
package main

import (
	"io"
	"log"
	"net/http"
)

func (h handler) getBody1() (string, error) {
	resp, err := h.client.Get(h.url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (h handler) getBody2() (string, error) {
	resp, err := h.client.Get(h.url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response: %v\n", err)
		}
	}()

	return string(body), nil
}

func (h handler) getStatusCode1(body io.Reader) (int, error) {
	resp, err := h.client.Post(h.url, "application/json", body)
	if err != nil {
		return 0, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response: %v\n", err)
		}
	}()

	return resp.StatusCode, nil
}

func (h handler) getStatusCode2(body io.Reader) (int, error) {
	resp, err := h.client.Post(h.url, "application/json", body)
	if err != nil {
		return 0, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response: %v\n", err)
		}
	}()

	_, _ = io.Copy(io.Discard, resp.Body)

	return resp.StatusCode, nil
}

type handler struct {
	client http.Client
	url    string
}
```

## What happens

### `getBody1` — body is never closed (resource leak)

`resp.Body` is read but never closed. The HTTP client's transport keeps the underlying connection open waiting for the body to be released, leaking file descriptors and exhausting the connection pool over time.

### `getBody2` — body is closed, but late

The body is closed with `defer`, which works because deferred calls run when the function returns. However, the `defer` is registered *after* `io.ReadAll`. The idiomatic and safer pattern is to close the body immediately after checking the error from `Get`, so the body is released even if reading fails.

### `getStatusCode1` — body is closed but not drained (connection can't be reused)

The response status is returned and the body is closed, but the body is never read or discarded. Because the body isn't fully consumed, Go's HTTP transport cannot return the connection to its keep-alive pool — the connection is closed instead of reused, hurting performance on repeated requests.

### `getStatusCode2` — correct

The body is closed with `defer` *and* fully drained with `io.Copy(io.Discard, resp.Body)`. This releases the connection back to the transport's pool so it can be reused by subsequent requests.

## The lesson

When making an HTTP request with `net/http`:

1. **Always close `resp.Body`** — even when you only care about the status code, and even on error paths. Register the `defer` right after checking the error from `Get`/`Post`, before reading the body.
2. **Always read or discard the body fully** — if you don't need the body, drain it with `io.Copy(io.Discard, resp.Body)` so the connection can be reused by the keep-alive pool.

The correct pattern:

```go
resp, err := client.Get(url)
if err != nil {
    return "", err
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
    return "", err
}
return string(body), nil
```

If you only need the status code:

```go
resp, err := client.Post(url, "application/json", body)
if err != nil {
    return 0, err
}
defer resp.Body.Close()

_, _ = io.Copy(io.Discard, resp.Body)

return resp.StatusCode, nil
```

## Run it

This directory is part of the root module (`module 100-go-mistakes`), so there's no local `go.mod`. Build and run from this directory:

```bash
go run .
```

## Related

Part of the [100 Go Mistakes](https://github.com/teivah/100-go-mistakes) series — mistake 79: closing resources.

- Sibling examples: `../os-file/` (file handles), `../sql-rows/` (database rows)
