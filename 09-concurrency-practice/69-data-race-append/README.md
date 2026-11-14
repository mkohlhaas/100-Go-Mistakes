# 69 - data race on append

Companion code for **chapter 9 (concurrency practice)**, exercise 69 of [*100 Go Mistakes and How to Avoid Them*](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them).

This program demonstrates that appending to the same slice from multiple goroutines is a **data race**. A slice is just a header (pointer, length, capacity) over a shared backing array, and concurrent appends can write to the same array position and race on the length field.

## What the code does

`main.go` provides two listings, `listing1` and `listing2`, run by `main()`.

### listing1 — appending to a shared slice (racy)

```go
func listing1() {
	s := make([]int, 1)

	go func() {
		s1 := append(s, 1)
		fmt.Println(s1)
	}()

	go func() {
		s2 := append(s, 1)
		fmt.Println(s2)
	}()
}
```

Both goroutines append to `s`. Because `s` has no spare capacity (`len == cap == 1`), each append happens to allocate a new backing array, so the race detector stays quiet — but the pattern is unsafe: give the slice any spare capacity (`make([]int, 0, 1)`) and both goroutines write to the same index of the **shared** backing array, which `go run -race` flags as a data race and can silently lose appends.

#### Possible outcomes of listing1

**As written** (`make([]int, 1)`):

- Both goroutines always print `[0 1]`.
- The order of the two lines is nondeterministic, but since both lines are identical, the output looks the same every run:

```
[0 1]
[0 1]
```

Each goroutine appends to its own freshly allocated backing array, so there are no shared writes and `go run -race` reports nothing. This is why the code looks safe even though the pattern is not.

**With spare capacity** (`make([]int, 0, 1)`):

- Both goroutines print `[1]`, but they write to the same index of the shared backing array concurrently — `go run -race` flags a data race:

```
==================
WARNING: DATA RACE
Write at 0x... by goroutine 11:
  main.listing1.func2()
      ...

Previous write at 0x... by goroutine 10:
  main.listing1.func1()
      ...

Found 1 data race(s)
exit status 66
```

- With a longer slice (e.g. `make([]int, 0, 2)` where both append past the current length), one goroutine can even overwrite the other's value, so appends can be silently lost.

### listing2 — copy before appending (safe)

```go
func listing2() {
	s := make([]int, 0, 1)

	go func() {
		sCopy := make([]int, len(s), cap(s))
		copy(sCopy, s)

		s1 := append(sCopy, 1)
		fmt.Println(s1)
	}()

	go func() {
		sCopy := make([]int, len(s), cap(s))
		copy(sCopy, s)

		s2 := append(sCopy, 1)
		fmt.Println(s2)
	}()
}
```

Each goroutine allocates a private backing array (`make` + `copy`) and appends to its own slice copy. No memory is shared between the goroutines, so there is no data race and each goroutine observes its own result (`[1]`).

## Key takeaway

- Copying a slice copies only the header — the pointer, length, and capacity — **not** the backing array. Two slices sharing a backing array that are appended to concurrently race on that array.
- Before appending from multiple goroutines, give each goroutine its own copy (as `listing2` does), serialize access with a `sync.Mutex`, or funnel values through a channel to a single appending goroutine.
- Always race-test concurrent code with `go run -race` / `go test -race`.

## How to run

`main()` runs both listings with a short sleep between them so the goroutines finish before the program exits:

```bash
go run .
```

Run with the race detector to catch the data race:

```bash
go run -race .
```

Example output:

```
=== listing1: appending to a shared slice ===
[0 1]
[0 1]
=== listing2: appending to private copies ===
[1]
[1]
done
```

## Observations

**listing1 with spare capacity** (`make([]int, 0, 1)`):

```
==================
WARNING: DATA RACE
Write at 0x... by goroutine 11:
  main.listing1.func2()
      ...

Previous write at 0x... by goroutine 10:
  main.listing1.func1()
      ...

==================
Found 1 data race(s)
exit status 66
```

**listing2** (private copy per goroutine): no race warning; both goroutines print their own slice:

```
[1]
[1]
```

## Requirements

- Go 1.27+

## Related

Part of the [100 Go Mistakes](https://github.com/teivah/100-go-mistakes) series — exercise 69: data race on slice append.
