# 21-slice-init — Inefficient slice initialization

This code demonstrates **mistake #21 from "100 Go Mistakes":** building a result slice by appending to a slice that starts with zero capacity, which forces repeated reallocations and copies.

## The mistake

```go
func convertEmptySlice(foos []Foo) []Bar {
    bars := make([]Bar, 0)   // length 0, capacity 0

    for _, foo := range foos {
        bars = append(bars, fooToBar(foo))
    }
    return bars
}
```

`make([]Bar, 0)` creates a slice with **zero length and zero capacity**. Every `append` that exceeds the current capacity forces Go to:

1. Allocate a new, larger backing array
2. Copy all existing elements over
3. Discard the old array

With one million input elements, that means many reallocations and copies — the most expensive way to build a slice.

## Better approaches

Pre-allocate so there is at most one allocation:

```go
// Option 1: pre-allocate capacity, append
func convertGivenCapacity(foos []Foo) []Bar {
    n := len(foos)
    bars := make([]Bar, 0, n)   // capacity known up front

    for _, foo := range foos {
        bars = append(bars, fooToBar(foo))
    }
    return bars
}

// Option 2: allocate exact length, assign by index (best when size is known)
func convertGivenLength(foos []Foo) []Bar {
    n := len(foos)
    bars := make([]Bar, n)

    for i, foo := range foos {
        bars[i] = fooToBar(foo)
    }
    return bars
}
```

- `make([]Bar, 0, n)` — reserves capacity up front, so `append` never reallocates. One allocation total.
- `make([]Bar, n)` — allocates once and assigns by index, avoiding `append` overhead entirely. Use this when you know the final size.

## How to run

From this directory:

```sh
go test -bench=. -benchmem -run=^$ .
```

- `-bench=.` runs all benchmarks
- `-benchmem` reports memory allocations per op
- `-run=^$ ` skips tests (there are none here anyway)

## Benchmark results

```
BenchmarkConvert_EmptySlice-16        456   2632533 ns/op   0 B/op   0 allocs/op
BenchmarkConvert_GivenCapacity-16     2202    543631 ns/op   0 B/op   0 allocs/op
BenchmarkConvert_GivenLength-16       2173    538403 ns/op   0 B/op   0 allocs/op
```

`convertEmptySlice` is roughly **5× slower** than the other two because of repeated reallocations and copying. The two pre-allocated versions are close, but `GivenLength` (assign by index) is slightly cleaner and avoids `append` overhead entirely.

## Takeaway

When you know how many elements a result slice will hold, don't start with `make([]T, 0)` and append. Either pre-allocate capacity with `make([]T, 0, n)`, or — ideally — allocate the exact length with `make([]T, n)` and populate by index.
