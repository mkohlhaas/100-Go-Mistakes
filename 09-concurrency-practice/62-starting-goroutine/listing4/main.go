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
