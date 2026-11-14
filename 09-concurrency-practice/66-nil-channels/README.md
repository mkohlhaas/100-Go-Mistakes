# 66 - Nil Channels

Demonstrates how **nil channels** can be used to merge multiple input channels into one output channel in Go, and why they matter for detecting when all inputs are closed.

## The problem

Merging two channels (`ch1`, `ch2`) into a single output channel (`ch`) sounds simple, but each of the first three attempts below has a flaw:

### `merge1` — sequential, blocking

```go
for v := range ch1 {
    ch <- v
}
for v := range ch2 {
    ch <- v
}
close(ch)
```

This drains `ch1` completely before touching `ch2`. If `ch1` is never closed (or is fed slowly while `ch2` has urgent data), `ch2` values stall. Worse, if `ch1` never closes, `ch2` is never read at all.

### `merge2` — select, but never terminates

```go
select {
case v := <-ch1:
    ch <- v
case v := <-ch2:
    ch <- v
}
```

A `select` picks a ready case, but when both input channels close, the `default`less `select` keeps blocking forever — there's no way to know both channels are done, so `close(ch)` is never reached and downstream `range` loops hang.

### `merge3` — select with closed flags

```go
case v, open := <-ch1:
    if !open {
        ch1Closed = true
        break
    }
```

This works, but it requires manual bookkeeping (`ch1Closed`, `ch2Closed`) and an extra termination check after every `select` iteration.

### `merge4` — the nil channel trick (idiomatic Go)

```go
for ch1 != nil || ch2 != nil {
    select {
    case v, open := <-ch1:
        if !open {
            ch1 = nil   // remove from select
            break
        }
        ch <- v
    case v, open := <-ch2:
        if !open {
            ch2 = nil   // remove from select
            break
        }
        ch <- v
    }
}
close(ch)
```

**Key insight:** receiving from a **nil channel blocks forever**, and a `select` case whose channel is nil is never selected. So once a channel is closed, setting it to `nil` removes it from the `select`. The loop exits when both channels are nil (both closed), and `close(ch)` signals the consumer that no more values are coming.

## Run it

`main.go` includes a `main()` entry point that demonstrates all four merges:

```bash
go run main.go
```

Output:

```
merge1: 1 2 3 4 
merge3: 1 2 3 4 
merge4: 1 2 3 4 
merge2: (never terminates — see the README for why)
```

`merge2` is intentionally left broken (its `select` has no way to detect that both inputs are closed, so `close(ch)` is never reached). `main()` runs `merge1`, `merge3`, and `merge4`, then prints a note about `merge2` instead of calling it — otherwise the program would hang forever.

## Channel read scenarios

Reading from a channel behaves differently depending on its state and the receive form used:

| Scenario | Channel status | Behavior |
|----------|----------------|----------|
| `v := <-ch` | Open | Blocks until a value is available, then assigns it to `v`. |
| `v, open := <-ch` | Open | Same as above; `open` is `true`. |
| `v := <-ch` | Closed, empty | Returns immediately with `v` set to the zero value (`0` for `int`) — no way to tell the channel is closed. |
| `v, open := <-ch` | Closed, empty | Returns immediately with `v` set to the zero value and `open` set to `false`. |
| `<-nil` | Nil | Blocks forever. A `select` case with a nil channel is never chosen. |
| `<-ch` | Closed, non-empty | Drains the remaining buffered values first, then behaves as "closed, drained" above. |

This is why `merge2` can't terminate: it uses `v := <-ch1`, so once both inputs are closed it keeps receiving `0` forever and never reaches `close(ch)`. `merge3` and `merge4` use `v, open := <-ch1` to detect closure (`open == false`), and `merge4` then sets the channel to `nil` so the `select` stops considering it.

**Buffering note:**

1. These scenarios assume **unbuffered** channels unless stated otherwise.
2. On an **unbuffered** channel, `<-ch` blocks until a sender is ready, and a closed channel is always empty immediately after closing.
3. On a **buffered** channel, `<-ch` returns immediately while the buffer has values, and a closed channel can still hold buffered values that must be drained first — that's the "closed, non-empty" row above.
4. None of the `merge*` functions use buffered input channels; the output channel `ch` is buffered with capacity 1.

## Key takeaways

- A **nil channel blocks forever** on send and receive.
- A `select` case with a **nil channel is effectively disabled** — it is never chosen.
- Setting a closed channel to `nil` is the idiomatic way to merge a dynamic set of channels without manual "closed" bookkeeping.
- Always `close` the output channel once all inputs are drained, so consumers using `range` can terminate.
