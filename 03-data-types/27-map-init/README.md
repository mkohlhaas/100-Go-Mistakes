# Map init — allocate maps with a size hint

This example shows how omitting the capacity hint when creating a map forces Go
to reallocate and rehash the map as it grows.

## The mistake

```go
// Slow: Go grows the map incrementally as keys are inserted.
m := make(map[int]struct{})
for j := range n {
    m[j] = struct{}{}
}
```

When you don't tell `make` how many entries to expect, the map starts small and
is resized several times as it fills up. Each resize allocates a new bucket
array and rehashes every existing key — wasteful when you know the final size
up front.

## The fix

```go
// Fast: reserve space for all n entries up front.
m := make(map[int]struct{}, n)
for j := range n {
    m[j] = struct{}{}
}
```

Passing `n` as the capacity hint allocates enough buckets for the whole load
once, avoiding reallocations during insertion.

## Benchmarks

Two benchmarks compare the two approaches with `n = 1_000_000`:

- `BenchmarkMapWithoutSize` — `make(map[int]struct{})` (no hint)
- `BenchmarkMapWithSize` — `make(map[int]struct{}, n)` (with hint)

## Running the benchmarks

Run both benchmarks with allocation stats:

```bash
go test -bench=. -benchmem
```

Sample output (numbers vary by machine):

```text
goos: linux
goarch: amd64
pkg: github.com/teivah/100-go-mistakes/03-data-types/27-map-init
cpu: AMD Ryzen 7 2700 Eight-Core Processor
BenchmarkMapWithoutSize-16          7   164777928 ns/op  75467379 B/op    8196 allocs/op
BenchmarkMapWithSize-16             9   115364514 ns/op  37832720 B/op    4097 allocs/op
PASS
ok  	github.com/teivah/100-go-mistakes/03-data-types/27-map-init	2.202s
```

The version with the size hint allocates about half as many times.

## Run a single benchmark

```bash
go test -bench=BenchmarkMapWithSize -benchmem
```

## Run with more iterations

```bash
go test -bench=. -benchtime=10x -benchmem
```

## Requirements

- Go 1.22+ (uses `b.Loop()` and the `1_000_000` literal syntax)
