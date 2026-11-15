# Mistake 76 — Using `time.After` in a loop

Chapter 10: Standard library · [`100 Go Mistakes and How to Avoid Them`](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them) · [100go.co](https://100go.co)

The consumers wait `warnAfter` for a message before warning — a package variable that defaults to `time.Hour` (the mistake as presented in the book). `main` overrides it to `2 * time.Second` so the demo fires quickly.

## The mistake

`consumer1` waits for events on a channel and logs a warning if no message arrives within `warnAfter` (defaults to an hour):

```go
select {
case event := <-ch:
    handle(event)
case <-time.After(warnAfter):
    log.Println("warning: no messages received")
}
```

`time.After(d)` starts a **new** `time.Timer` under the hood and returns only its channel. Because the select sits inside a `for` loop, every iteration creates another timer. None of these timers can be stopped — `time.After` does not expose the `*time.Timer` — so they pile up and are only released by the garbage collector after they fire an hour later. That is a timer leak.

## The fix

`consumer2` bounds the wait with a context and **cancels** it as soon as an event is handled:

```go
ctx, cancel := context.WithTimeout(context.Background(), warnAfter)
select {
case event := <-ch:
    cancel()
    handle(event)
case <-ctx.Done():
    log.Println("warning: no messages received")
}
```

Canceling the context releases its timer immediately, so nothing leaks. Use `context.WithTimeout` (or `context.WithDeadline`) whenever you need a timeout around a channel operation.

One detail worth noting: `cancel()` cancels **only the current iteration's timer**, not the surrounding `for` loop. After handling an event, the loop immediately starts a new iteration with a fresh context and waits for the next message. That is why `consumer2` handles both events the producer sends — once per loop iteration — rather than exiting after the first. If you want the consumer to stop after a single event, you need to break out of the loop explicitly (e.g. with a `return` or `break`) in addition to canceling.

Is `cancel()` useless then? **No.** Without it, the context's timer would keep running for the full `warnAfter` duration after every handled event, firing pointlessly into an unobserved `ctx.Done()` channel before being garbage collected. Calling `cancel()` stops that timer immediately, releasing its resources right away. It's good hygiene — always cancel a cancelable context when you're done with it — even though in this loop structure it doesn't change which events get handled.

## About consumer3

`consumer3` reuses a single timer with `time.NewTimer` and `timer.Reset`, which avoids creating a new timer every iteration. The timer duration is `warnAfter`. However, calling `timer.Reset` on a timer that has already expired is racy: if a channel event arrives at almost the same instant as the timer firing, the timer can be reset and fire again almost immediately, producing a spurious warning. When in doubt, prefer `context.WithTimeout` (consumer2).

## Running

```sh
go run .
```

The demo runs each consumer against its own channel for a few seconds. To keep the demo short, `main` overrides the default hour-long timeout with `2 * time.Second`, so you can see the "warning: no messages received" fire after the producer stops sending. The Go module is resolved from the repository root (`go.mod`), so no `go.mod` is needed in this directory.

In production code the timeout would be `time.Hour` (or another real deadline), so the warning only fires after that long of silence on the channel.