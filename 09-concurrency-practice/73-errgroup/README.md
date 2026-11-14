# 73 — errgroup

Compares two ways of running work concurrently and collecting errors:

1. **`handler1`** — manual `sync.WaitGroup` (error-prone)
2. **`handler2`** — `golang.org/x/sync/errgroup` (recommended)

## The problem

When fanning out work to multiple goroutines, you need to:

* wait for all goroutines to finish,
* propagate the first error back to the caller,
* avoid data races on shared state (e.g. a results slice).

With a plain `sync.WaitGroup` there is no built-in error propagation, so you must add your own locking and error tracking (see `handler1`).

## The solution

`errgroup.WithContext` wraps a `context.Context` and a `sync.WaitGroup`-like group:

* `g.Go(f)` runs `f` in a goroutine; `f` returns an `error`.
* `g.Wait()` blocks until all goroutines finish and returns the **first non-nil error**.
* If any goroutine returns an error, the derived context is canceled, signaling the other goroutines to stop early.

## Files

* `main.go` — both implementations plus a `main()` that runs them.
* `README.md` — this file.

## Run it

```bash
go run .
```

Expected output:

```
handler1 error: circle "broken" has non-positive radius -1
```

The `"broken"` circle has a negative radius, so `foo` returns an error and both handlers propagate it. `handler2` exits early without printing results.

## Structs

```go
type Point struct {
    X, Y float64
}

type Circle struct {
    Name   string
    Center Point
    Radius float64
}

type Result struct {
    Name string
    Area float64
}
```

`computeCircleArea` validates that `Radius` is positive and computes the area (`πr²`).

## Key takeaway

Prefer `golang.org/x/sync/errgroup` over rolling your own `sync.WaitGroup` + error-handling logic. It gives you error propagation, context cancellation, and race-safe synchronization in a small, well-tested package.
