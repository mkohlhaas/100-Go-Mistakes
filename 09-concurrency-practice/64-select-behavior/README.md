# 64 - select behavior

Companion code for **chapter 9 (concurrency practice)**, exercise 64 of [*100 Go Mistakes and How to Avoid Them*](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them).

This program demonstrates how Go's `select` statement behaves when multiple channels are ready at the same time, and why you may need to **drain** a channel before exiting on a shutdown signal.

## What the code does

`main` spawns a goroutine running a `select` loop over two channels:

- `messageCh` — a buffered channel of `int` (capacity 10)
- `disconnectCh` — a signal channel of `struct{}`

Main sends the values `0` through `9` into `messageCh`, then signals `disconnectCh` and sleeps briefly to let the goroutine finish.

Two listings are provided:

### listing1 — exit immediately on disconnect

```go
func listing1(messageCh <-chan int, disconnectCh chan struct{}) {
	for {
		select {
		case v := <-messageCh:
			fmt.Println(v)
		case <-disconnectCh:
			fmt.Println("disconnection, return")
			return
		}
	}
}
```

When `disconnectCh` becomes ready while `messageCh` still holds values, both cases are ready. `select` picks **randomly** among the ready cases, so the goroutine may return on the disconnect before draining the remaining messages. Some values can be lost.

### listing2 — drain messages before disconnecting

```go
func listing2(messageCh <-chan int, disconnectCh chan struct{}) {
	for {
		select {
		case v := <-messageCh:
			fmt.Println(v)
		case <-disconnectCh:
			for {
				select {
				case v := <-messageCh:
					fmt.Println(v)
				default:
					fmt.Println("disconnection, return")
					return
				}
			}
		}
	}
}
```

When `disconnectCh` fires, the goroutine enters a nested loop that drains `messageCh` using non-blocking receives (`default` case) until the channel is empty, and only then returns. All buffered values are processed before shutdown.

### listing3 — drain before disconnecting, simplified

```go
func listing3(messageCh <-chan int, disconnectCh chan struct{}) {
	// Phase 1: process messages until disconnect is signaled.
OUTER:
	for {
		select {
		case v := <-messageCh:
			fmt.Println(v)
		case <-disconnectCh:
			break OUTER
		}
	}

	// Phase 2: drain remaining buffered messages before shutting down.
	for {
		select {
		case v := <-messageCh:
			fmt.Println(v)
		default:
			fmt.Println("disconnection, return")
			return
		}
	}
}
```

`listing3` has the same behavior as `listing2` — every buffered value is processed before the goroutine returns — but the code is simpler: a single `select` loop runs until disconnect, then a flat non-blocking drain loop finishes the remaining messages. There's no nested `select` inside a loop inside a `select`.

## Key takeaway

A `select` with multiple ready cases chooses one **uniformly at random**. If your goroutine must finish processing pending work on shutdown, don't return on the first disconnect signal — drain the channel first (as in `listing2`, or its simplified equivalent `listing3`), or you risk dropping messages.

## How to run

From this directory:

```bash
go run main.go
```

`main` runs all listings twice — once with a buffered `messageCh` (capacity 10) and once with an unbuffered `messageCh`:

1. `listing1` — exits immediately on disconnect (messages may be dropped)
2. `listing2` — drains all buffered messages before disconnecting
3. `listing3` — same as `listing2`, simplified (no nested `select`)

## Observations

**Buffered channel.** Because of random case selection in `listing1`, running the program a few times shows different subsets of `0`–`9` printed before `"disconnection, return"`, and sometimes nothing at all — values still sitting in the buffer can be lost. With `listing2` and `listing3`, all 10 values are always printed before the disconnection message because the goroutine drains the buffer first.

**Unbuffered channel.** All three listings print all 10 values in order (`0` through `9`) before disconnecting. An unbuffered channel is a rendezvous point: `main`'s send blocks until the goroutine receives, and the goroutine's receive blocks until `main` sends — neither side can proceed without the other. Since nothing can ever be buffered, nothing can be dropped, so the difference between `listing1` and `listing2`/`listing3` effectively disappears.

Example output:

```
$ go run main.go
=== buffered channel (cap 10) ===
--- listing1: exit immediately on disconnect ---
0
1
disconnection, return

--- listing2: drain messages before disconnecting ---
0
1
2
3
4
5
6
7
8
9
disconnection, return

--- listing3: drain before disconnecting, simplified ---
0
1
2
3
4
5
6
7
8
9
disconnection, return

=== unbuffered channel ===
--- listing1: exit immediately on disconnect ---
0
1
2
3
4
5
6
7
8
9
disconnection, return

--- listing2: drain messages before disconnecting ---
0
1
2
3
4
5
6
7
8
9
disconnection, return

--- listing3: drain before disconnecting, simplified ---
0
1
2
3
4
5
6
7
8
9
disconnection, return
```

## Requirements

- Go 1.27+
