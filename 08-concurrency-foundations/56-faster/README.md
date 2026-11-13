# 56-faster — Mergesort benchmark

This directory contains three implementations of mergesort side-by-side, plus Go
benchmarks to compare them:

| Implementation | Description |
|---|---|
| `sequentialMergesort` | Straightforward recursive mergesort |
| `parallelMergesortV1` | Spawns two goroutines per recursive call (naive parallelization) |
| `parallelMergesortV2` | Falls back to sequential mergesort below a `max` size threshold |

The goal of this example is to show that naive concurrency can be *slower* than
a sequential solution due to goroutine scheduling overhead, and how a threshold
cutoff (V2) mitigates that.

## Prerequisites

- Go 1.24 or newer (benchmarks use `testing.B.Loop()`, introduced in Go 1.24)

Check your version:

```sh
go version
```

## Running the benchmarks

From this directory:

```sh
go test -bench=. -benchmem -count=3
```

### Explanation of flags

| Flag | Meaning |
|---|---|
| `-bench=.` | Run all benchmarks matching the regex (`.` matches everything) |
| `-benchmem` | Report memory allocations per operation |
| `-count=3` | Run each benchmark 3 times to smooth out noise |
| `-run=^$` | Skip tests (no tests exist here; avoids `--- PASS` clutter) |

### Run a single benchmark

```sh
go test -run=^$ -bench=Benchmark_sequentialMergesort -benchmem -count=5
```

### Run all three and compare

```sh
go test -run=^$ -bench=. -benchmem -count=5
```

### CPU profile

To see where time is spent (useful for understanding goroutine overhead):

```sh
go test -run=^$ -bench=Benchmark_parallelMergesortV1 -benchmem -cpuprofile=cpu.prof
go tool pprof -top cpu.prof
```

## Benchmark setup

Each benchmark sorts a fresh slice of 10,000 random integers per iteration.
Timing stops around slice generation so only the sort itself is measured:

```go
for b.Loop() {
    b.StopTimer()
    input := getRandomElements()
    b.StartTimer()
    // sort
}
```

The sorted result is assigned to a package-level `global` variable so the
compiler cannot elide the sort call.

## Expected outcome

V1 typically runs *slower* than the sequential version because every recursive
split launches two goroutines, and the scheduler overhead dwarfs the parallel
work at small slice sizes. V2 usually beats V1 because it switches to the
sequential implementation once `len(s) <= max` (currently `max = 2048`),
avoiding excessive goroutine creation.

Run the benchmarks and compare `ns/op` across the three implementations to see
the effect yourself.