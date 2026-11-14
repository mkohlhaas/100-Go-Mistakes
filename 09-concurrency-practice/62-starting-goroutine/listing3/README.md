# Listing 3 — Returning the watcher and closing it

This example is part of **mistake #62** ("Starting goroutines without considering
their lifecycle") from *100 Go Mistakes* (chapter 09 — Concurrency practice).

## What this example shows

```go
package main

func main() {
	w := newWatcher()
	defer w.close()

	// Run the application
}

func newWatcher() watcher {
	w := watcher{}
	go w.watch()
	return w
}

type watcher struct { /* Some resources */ }

func (w watcher) watch() {}

func (w watcher) close() {
	// Close the resources
}
```

Unlike listing 1, `newWatcher()` returns the `watcher` instead of discarding it,
and `main` defers `w.close()` so the watcher's resources are released when
`main` returns. The caller — not the constructor — owns the watcher's
lifecycle.

## How it relates to the other listings

| Listing | Approach | Status |
|---------|----------|--------|
| `listing1/` | Goroutine started inside the constructor; nothing returned; no cleanup | ❌ The mistake — the watcher can never be stopped or closed |
| `listing2/` | `context.Context` passed into the goroutine for cancellation | ⚠️ Cancellation works, but the watcher is not returned and its resources have no explicit `close` |
| `listing3/` (this directory) | Watcher returned from the constructor; closed with `defer w.close()` | ⚠️ Ownership + resource cleanup, but no way to stop the goroutine |
| `listing4/` | Returned watcher + context cancellation + `close()` waits for the goroutine to exit | ✅ Lifecycle fully managed |

Listing 3 is a step in the right direction: the constructor hands the watcher
back to its caller, and `defer w.close()` releases the watcher's resources.
But it still starts the goroutine inside `newWatcher()` with no stop signal,
so there is no way to tell that goroutine to exit before resources are closed.
See `listing4/` for the complete fix that adds context cancellation and waits
for the goroutine to finish.

## Run it

```bash
go run main.go
```

The program does no observable work — it is a skeleton illustrating the
goroutine-lifecycle pattern.

## Key takeaways

- Don't start a goroutine inside a constructor and throw the object away;
  return the object so the caller controls its lifecycle.
- Pair the returned object with a `Close()` method and `defer close()` so
  resources are released deterministically when the caller is done.
- Combine an explicit stop signal (context cancellation, as in listing 2) with
  explicit resource cleanup (as here) for a fully managed lifecycle.
