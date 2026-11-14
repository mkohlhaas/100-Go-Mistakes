# Listing 4 — Cancel the goroutine, then close the watcher

⚠️ **Note:** This listing was created as part of this exercise and is **not part of the book code** from *100 Go Mistakes*. It is an extension of the book's listings 1–3 that completes the goroutine-lifecycle fix.

This example illustrates the idea behind **mistake #62** ("Starting goroutines without considering
their lifecycle") from *100 Go Mistakes* (chapter 09 — Concurrency practice).

## What this example shows

```go
package main

import (
	"context"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := newWatcher(ctx)
	defer w.close()

	// Run the application
}

type watcher struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	/* Some resources */
}

func newWatcher(ctx context.Context) *watcher {
	ctx, cancel := context.WithCancel(ctx)
	w := &watcher{
		ctx:    ctx,
		cancel: cancel,
	}
	w.wg.Add(1)
	go w.watch()
	return w
}

func (w *watcher) watch() {
	defer w.wg.Done()

	// Do work until the context is cancelled.
	<-w.ctx.Done()
}

func (w *watcher) close() {
	// 1. Ask the goroutine to stop.
	w.cancel()
	// 2. Wait for it to actually exit before releasing resources.
	w.wg.Wait()
	// 3. Close the resources.
}
```

Listing 4 combines the two ideas from the earlier listings:

- **listing 3** — `newWatcher` returns the watcher, so the caller owns it and
  defers `w.close()`.
- **listing 2** — a `context.Context` flows into the goroutine so it can be
  told to stop.

`close()` performs an ordered shutdown:

1. `w.cancel()` signals the goroutine via `ctx.Done()`.
2. `w.wg.Wait()` blocks until the goroutine has actually exited.
3. The watcher's resources are closed, safe in the knowledge that no goroutine
   is still using them.

## How it relates to the other listings

| Listing | Approach | Status |
|---------|----------|--------|
| `listing1/` | Goroutine started inside the constructor; nothing returned; no cleanup | ❌ The mistake — the watcher can never be stopped or closed |
| `listing2/` | `context.Context` passed into the goroutine for cancellation | ⚠️ Cancellation works, but the watcher is not returned and its resources have no explicit `close` |
| `listing3/` | Watcher returned from the constructor; closed with `defer w.close()` | ⚠️ Ownership + resource cleanup, but no way to stop the goroutine |
| `listing4/` (this directory) | Returned watcher + context cancellation + `close()` waits for exit | ✅ Lifecycle fully managed |

Listing 4 is the complete, idiomatic fix: the caller owns the watcher, can
cancel its work, waits for the goroutine to finish, and then releases
resources — all deterministically via `defer w.close()`.

## Run it

```bash
go run main.go
```

The program does no observable work — it is a skeleton illustrating the
goroutine-lifecycle pattern.

## Key takeaways

- Return the object that owns the goroutine; let the caller decide its
  lifetime with `defer close()`.
- Pass a `context.Context` into long-running goroutines so cancellation can
  reach them.
- In `Close()`, cancel first, `Wait()` for the goroutine to exit, then release
  resources — never free resources a goroutine might still be using.
