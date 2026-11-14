# Listing 1 — Starting a goroutine inside a constructor

Part of **100 Go Mistakes**, Chapter 09: Concurrency Practice, Example 62.

## What this example shows

```go
func main() {
	newWatcher()

	// Run the application
}

func newWatcher() {
	w := watcher{}
	go w.watch()
}

type watcher struct { /* Some resources */ }

func (w watcher) watch() {}
```

`newWatcher()` creates a `watcher` and immediately launches `w.watch()` in its own goroutine before returning to the caller.

## The problem

Starting a goroutine inside a constructor (or factory function) removes the caller's control over the goroutine's lifecycle:

1. **No control over when work starts.** The goroutine begins running before `newWatcher` returns, so the caller can't finish setup (e.g. registering handlers, opening channels) before work begins.
2. **Value receiver = copy semantics.** `watch()` has a value receiver (`func (w watcher) watch()`), so the goroutine operates on a *copy* of the watcher. Any state the goroutine mutates is invisible to the caller's instance.
3. **No way to shut down.** There is no `done` channel, `context.Context`, or `Close()` method, so the goroutine can't be stopped and its resources can't be released when the application exits.

## How to fix it

Return the watcher from the constructor and let the caller decide when (and whether) to start the goroutine:

```go
func newWatcher() *watcher {
	return &watcher{}
}

func (w *watcher) watch(ctx context.Context) { ... }

func main() {
	w := newWatcher()
	defer w.close()

	// caller decides when to start work
	go w.watch(context.Background())

	// Run the application
}
```

See the sibling listings in this directory for incremental fixes:
- **listing2** — threads a `context.Context` through so the goroutine can be cancelled.
- **listing3** — returns the watcher and adds a `close()` method so the caller owns the lifecycle.
- **listing4** — combines context cancellation with the returned watcher so `close()` can stop the goroutine and wait for it to exit before releasing resources.

## Running

```bash
go run main.go
```

This listing compiles but does nothing observable — it is intended to illustrate the pattern, not to be run as a standalone application.