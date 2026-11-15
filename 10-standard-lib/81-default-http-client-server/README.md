# Mistake 81 — Default HTTP client and server without timeouts

Go's `net/http` package does not apply timeouts by default. A request (client side) or a
connection (server side) that never times out can hang forever, leaking goroutines, file
descriptors, and memory until the process runs out of resources.

This directory shows the **correct** way to configure timeouts on both sides of an HTTP
exchange.

## The problem

Without explicit timeouts:

* **Client:** `http.Get`, `http.Post`, `http.Client.Do` block indefinitely if the server
  never responds or responds slowly.
* **Server:** `http.ListenAndServe` keeps each connection alive indefinitely. A slow or
  malicious client can occupy a goroutine and a connection forever.

## The fix

Set timeouts explicitly on `http.Client` / `http.Transport` and on `http.Server`.

### Client — `client/main.go`

```go
client := &http.Client{
    Timeout: 5 * time.Second,
    Transport: &http.Transport{
        DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
        TLSHandshakeTimeout:   time.Second,
        ResponseHeaderTimeout: time.Second,
    },
}
```

* `http.Client.Timeout` — overall limit for the request (including body read).
* `http.Transport` timeouts — per-connection limits: dial, TLS handshake, and waiting for
  response headers.

### Server — `server/main.go`

```go
s := &http.Server{
    Addr:              ":8080",
    ReadHeaderTimeout: 500 * time.Millisecond,
    ReadTimeout:       500 * time.Millisecond,
    Handler:           http.TimeoutHandler(handler{}, time.Second, "foo"),
}
```

* `http.Server.ReadHeaderTimeout` — limit for reading request headers.
* `http.Server.ReadTimeout` — limit for reading the full request body.
* `http.TimeoutHandler` — wraps each handler with a per-request timeout.

## Project structure

```
81-default-http-client-server/
├── client/
│   └── main.go   # HTTP client with explicit timeouts
├── server/
│   └── main.go   # HTTP server with explicit timeouts
└── README.md
```

## Running

From this directory:

```bash
# client (compiles only — the client is not wired to a live server here)
cd client && go run .

# server (compiles only — the server is not started here)
cd server && go run .
```

These snippets illustrate the configuration; wire them into your own `http.Client` and
`http.Server` (and call `s.ListenAndServe()` / `client.Do(req)`) to use them.

## Related

* [net/http — Package http](https://pkg.go.dev/net/http)
* [http.Client](https://pkg.go.dev/net/http#Client)
* [http.Transport](https://pkg.go.dev/net/http#Transport)
* [http.Server](https://pkg.go.dev/net/http#Server)
* [http.TimeoutHandler](https://pkg.go.dev/net/http#TimeoutHandler)
