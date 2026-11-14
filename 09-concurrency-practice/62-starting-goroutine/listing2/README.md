# 62 — Starting goroutines (listing 2)

This example is part of **mistake #62** ("Starting goroutines without considering
their lifecycle") from *100 Go Mistakes* (chapter 09 — Concurrency practice).

## The mistake

A goroutine that is started with `go` but never given a way to stop leaks for the
entire lifetime of the program:

```go
func newWatcher() {
	w := watcher{}
	go w.watch() // nobody can tell this goroutine to stop
}
```

If `watcher` holds resources (file handles, DB connections, timers, …) those
resources are never released, and `main` cannot shut the watcher down cleanly.

## What listing 2 shows

Listing 2 threads a `context.Context` into the watcher so cancellation can
signal the goroutine to stop:

```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newWatcher(ctx)

	// Run the application
}

func newWatcher(ctx context.Context) {
	w := watcher{}
	go w.watch(ctx)
}

func (w watcher) watch(context.Context) {}
```

`context.WithCancel` + `defer cancel()` gives `main` a handle to request
shutdown, and the watcher's `watch` method can observe `ctx.Done()` and exit.

## How it relates to the other listings

| Listing | Approach | Status |
|---------|----------|--------|
| `listing1/` | Goroutine started with no context, no cleanup | ❌ The mistake — the watcher can never be stopped |
| `listing2/` (this directory) | `context.Context` passed into the goroutine | ⚠️ Cancellation works, but the watcher is not returned and its resources have no explicit `close` |
| `listing3/` | Watcher is returned and closed with `defer w.close()` | ⚠️ Ownership + resource cleanup, but no way to stop the goroutine |
| `listing4/` | Returned watcher + context cancellation + `close()` waits for the goroutine to exit | ✅ Lifecycle fully managed |

Listing 2 is a step in the right direction (cancellation via context), but
because `newWatcher` returns nothing, callers cannot explicitly close the
watcher's resources. See `listing3/` for ownership + cleanup, and `listing4/`
for the complete, idiomatic fix that also stops the goroutine before closing.

## Run it

```bash
go run main.go
```

The program does no observable work — it is a skeleton illustrating the
goroutine-lifecycle pattern.

## Key takeaways

- Never start a goroutine you cannot stop.
- Pass a `context.Context` into long-running goroutines so cancellation can
  reach them.
- Prefer returning the object that owns the goroutine and exposing a `Close`
  method, used with `defer`, so both cancellation and resource cleanup are
  handled by the caller.