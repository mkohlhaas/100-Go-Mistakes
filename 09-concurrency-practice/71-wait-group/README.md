# WaitGroup examples

This package demonstrates correct and incorrect usage of `sync.WaitGroup` in Go,
based on practice 71 (WaitGroup) from *100 Go Mistakes and How to Avoid Them*.

The file `main.go` contains three listings, each launching three goroutines that
increment a shared counter and print the final value.

## Listing 1 — ❌ Mistake: calling `Add` inside the goroutine

```go
wg := sync.WaitGroup{}
var v uint64

for range 3 {
    go func() {
        wg.Add(1)          // called after the goroutine starts
        atomic.AddUint64(&v, 1)
        wg.Done()
    }()
}

wg.Wait()                  // may return before any Add happens
fmt.Println(v)
```

**The problem.** `WaitGroup.Add` and `Wait` must not run concurrently. Because
the goroutine may not have started (and thus may not have called `Add`) by the
time the main goroutine reaches `wg.Wait()`, `Wait` can return with an internal
counter of `0` — or even a negative counter, which panics. The printed result
is therefore unreliable.

## Listing 2 — ✅ Pre-adding the total count

```go
wg := sync.WaitGroup{}
var v uint64

wg.Add(3)                  // register all work up front
for range 3 {
    go func() {
        atomic.AddUint64(&v, 1)
        wg.Done()
    }()
}

wg.Wait()
fmt.Println(v)
```

**Why this works.** The total count is registered *before* any goroutine
starts, so there is no race between `Add` and `Wait`. Each goroutine signals
completion with `Done`.

Use this form when you know the number of tasks up front.

## Listing 3 — ✅ Adding per task with `wg.Go`

```go
wg := sync.WaitGroup{}
var v uint64

for range 3 {
    wg.Go(func() {
        atomic.AddUint64(&v, 1)
    })
}

wg.Wait()
fmt.Println(v)
```

**Why this works.** `wg.Go` (added in Go 1.24) internally calls `Add(1)` and
launches the goroutine atomically, so the counter is always registered before
`Wait` can observe it. This is the cleanest form when you know the number of
tasks up front, and it removes the risk of forgetting `Done` — `wg.Go` handles
it for you.

Equivalent manual form (for older Go versions):

```go
for range 3 {
    wg.Add(1)              // register work before the goroutine starts
    go func() {
        defer wg.Done()
        atomic.AddUint64(&v, 1)
    }()
}

wg.Wait()
fmt.Println(v)
```

`wg.Add(1)` happens-before the goroutine is created, so the counter is
incremented before `Wait` can observe it. This is safe even when the number of
tasks isn't known until runtime.

## Key rules for `sync.WaitGroup`

1. **Call `Add` before `Wait`.** `Add` and `Wait` must not run concurrently.
2. **Call `Add` before starting the goroutine.** Otherwise the goroutine may
   increment the counter after `Wait` has already returned.
3. **Call `Done` exactly once per `Add(1)`**, typically with `defer wg.Done()`
   as the first statement in the goroutine.
4. **Prefer `wg.Go` (Go 1.24+)** which combines `Add(1)` and `go` atomically and
   calls `Done` for you, eliminating this entire class of mistakes.
5. **Never copy a `WaitGroup`** after its first use — pass it by pointer or keep
   it in one place.

## Running

```bash
go run main.go
```

Expected output:

```
listing1:
0

listing2:
3

listing3:
3
```

Listing 1 prints `0` (or another value below 3, or panics) because `Wait`
returns before the goroutines call `Add`. Listings 2 and 3 reliably print `3`.
