# Mistake #39 — Building Strings by Repeated `+=`

This code demonstrates why concatenating strings inside a loop with `s += value` is slow (and how `strings.Builder` fixes it).

## The code

Three ways to join `[]string`:

* `concat1` — `s += value` in a loop (the mistake)
* `concat2` — `strings.Builder`
* `concat3` — `strings.Builder` with `sb.Grow(total)` pre-allocation

### Why `concat1` is slow

Go strings are immutable, so every `s += value` allocates a new backing array and copies the whole string built so far. With `n` values of length `k`, that's O(n²·k) work — copies grow quadratically as the string gets longer.

`strings.Builder` accumulates the pieces in an internal buffer with a single allocation growth strategy, so the same join is O(n·k). `concat3` goes one step further by computing the total length up front and calling `sb.Grow(total)`, avoiding intermediate buffer reallocations entirely.

## Run the benchmarks

Go 1.24+ is required (the benchmarks use `b.Loop()`).

From this directory:

```bash
go test -bench=. -benchmem -run=^$ .
```

That runs all three benchmarks (`BenchmarkConcatV1`, `BenchmarkConcatV2`, `BenchmarkConcatV3`) and prints timings plus per-op memory allocation:

```
goos: linux
goarch: amd64
pkg: main
cpu: ...
BenchmarkConcatV1-16     100000     12345 ns/op    1000000 B/op    1000 allocs/op
BenchmarkConcatV2-16     500000      3456 ns/op      1000000 B/op      10 allocs/op
BenchmarkConcatV3-16     500000      3000 ns/op      1000000 B/op       1 allocs/op
PASS
ok      main    1.234s
```

You should see `V1` far slower than `V2`, and `V2` slightly slower (or equal) to `V3`, with `V1` allocating orders of magnitude more.

### Run a single benchmark

```bash
go test -bench=BenchmarkConcatV1 -benchmem -run=^$ .
```

### Run with CPU/memory profiling

```bash
go test -bench=. -benchmem -run=^$ -cpuprofile=cpu.prof -memprofile=mem.prof .
go tool pprof -top cpu.prof
go tool pprof -top mem.prof
```

## Key takeaway

Never build a string by repeatedly appending with `+` or `+=` in a loop. Use `strings.Builder`, and call `sb.Grow` with the expected total size when you know it in advance to skip reallocations entirely.
