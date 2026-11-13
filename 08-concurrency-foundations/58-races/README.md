# Memory model

Examples from [*100 Go Mistakes and How to Avoid Them*](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them), mistake 58 — races and the Go memory model.

## What this shows

Go does not guarantee that operations in different goroutines are observed in program order. For one goroutine to see another's writes, a **happens-before** edge must exist between them. Without it, you have a data race: the result is nondeterministic, and `go run -race` will report it.

## The memory-model package (`memory-model/`)

| Listing | Synchronized? | Typical output | Notes |
|---------|---------------|----------------|-------|
| [listing1](main.go#L5) | No | (none) | Goroutine writes `i`, but nothing observes it and nothing orders the write. |
| [listing2](main.go#L12) | No | `0` or `1` | Classic data race: the goroutine increments `i` while `main` reads it. Flagged by `-race`. |
| [listing3](main.go#L20) | Yes | `1` | Channel send/receive rendezvous creates a happens-before edge: `i++` happens before the print. |
| [listing4](main.go#L31) | No | `1`, or nothing | `close(ch)` unblocks the receiver, and per the Go memory model **close does synchronize with the receive** that returns a zero value: `i++` happens before `close(ch)`, which happens before `<-ch`, which happens before `fmt.Println(i)`. So `i` is guaranteed to be `1` if it prints at all. Also, `main()` returns right after `close(ch)`, so the program can exit before the goroutine runs `fmt.Println(i)` — printing nothing. |
| [listing5](main.go#L42) | No | `0` or `1` | Buffered send succeeds immediately without a receiver, so it establishes no ordering with the goroutine. Flagged by `-race`. |
| [listing6](main.go#L53) | Yes | `1` | Unbuffered channel forces a rendezvous: the send completes only when the goroutine receives, after it has written `i = 1`. |

## How to run

```bash
cd 08-concurrency-foundations/58-races/memory-model
go run .
go run -race .
```

Run `go run .` several times for listings 2, 4 and 5 — their output can vary between runs. With `-race`, listings 2 and 5 are reported as data races.

## The races package (`races/`)

The `races` directory shows a single race — two goroutines incrementing the
same counter — fixed four ways:

| Listing | Fix | Safe? | Typical output | Notes |
|---------|-----|-------|----------------|-------|
| [listing1](races/main.go#L9) | None | No | 0, 1 or 2 | Unsynchronized `i++`: a data race, flagged by `go run -race`; the result is nondeterministic. |
| [listing2](races/main.go#L27) | `sync/atomic` | Yes | 2 | `atomic.AddInt64` makes the increments race-free. |
| [listing3](races/main.go#L45) | `sync.Mutex` | Yes | 2 | The mutex serializes access to `i`. |
| [listing4](races/main.go#L68) | Unbuffered channel | Yes | 2 | The sends rendezvous with the receives, so both increments are observed. |
| [listing5](races/main.go#L85) | `sync.Mutex` | Yes | 1 or 2 | Both goroutines write (`i = 1` and `i = 2`) under the mutex; the last write wins, so the result is nondeterministic but race-free. |

### How to run

```bash
cd 08-concurrency-foundations/58-races/races
go run .
go run -race .
```

`listing1` is reported as a data race by `-race`; listings 2–5 are clean. Run
listing 1 several times — its result varies between runs.

## Rules of thumb

1. **Channel send/receive synchronizes.** A send happens before the corresponding receive completes (and vice versa for the rendezvous on an unbuffered channel), creating a happens-before edge between the two goroutines.
2. **`close` synchronizes with the receive that returns a zero value.** Per the memory model, closing a channel happens before a receive that returns the zero value because the channel is closed and empty. So `close` *does* create a happens-before edge in that case — but it only synchronizes the channel operation itself, not arbitrary memory writes unless you establish the chain through program order as well.
3. **A non-blocking buffered send does not synchronize.** If the send completes because there is buffer space (not because a receiver took the value), no happens-before edge is created.
4. **When in doubt, use `sync.Mutex` or `sync/atomic`**, and verify with `go test -race`. Never rely on thread scheduling or `time.Sleep` for correctness.
