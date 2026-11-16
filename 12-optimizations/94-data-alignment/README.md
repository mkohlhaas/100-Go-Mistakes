# Data Alignment

This example demonstrates how struct field ordering affects memory layout and performance in Go.

## The Mistake

`Foo1` orders its fields as `byte`, `int64`, `byte`. Because `int64` requires 8-byte alignment, the compiler inserts padding bytes between `b1` and `i`, and after `b2`, so each instance occupies 24 bytes with gaps:

```
Foo1: [b1][pad][pad][pad][pad][pad][pad][pad][i ........][b2][pad][pad][pad][pad][pad][pad][pad]
      1   7 padding bytes        8 bytes         1   7 padding bytes
```

That wastes space and increases memory bandwidth usage when iterating over a slice.

## The Fix

`Foo2` orders its fields as `int64`, `byte`, `byte`, so the `int64` is naturally aligned and only 6 trailing padding bytes are needed:

```
Foo2: [i ........][b1][b2][pad][pad][pad][pad][pad][pad]
      8 bytes     1   1   6 padding bytes
```

## Result

`Foo2` is 16 bytes per element while `Foo1` is 24 bytes, so summing over a slice of `Foo2` reads less memory per element and is more cache-friendly than summing over `Foo1`.

Run the benchmarks:

```bash
go test -bench=. -benchmem
```

### Example output

```
goos: linux
goarch: amd64
pkg: 100-go-mistakes/12-optimizations/94-data-alignment
cpu: Intel(R) Core(TM) i5-5250U CPU @ 1.60GHz
BenchmarkSum1-4   	   1000	   1047894 ns/op	       0 B/op	       0 allocs/op
BenchmarkSum2-4   	    760	   1571569 ns/op	       0 B/op	       0 allocs/op
```

Both benchmarks report `0 B/op` and `0 allocs/op` because the slices are allocated once outside the loop. Absolute timings vary between runs and across machines (on this CPU they're noisy and not consistently in favor of either), so timing alone isn't the reliable signal here. The deterministic win is the memory layout: `Foo2` uses 16 bytes per element versus `Foo1`'s 24 bytes (33% less), which reduces memory bandwidth pressure and improves cache utilization at scale.

## Rule of Thumb

Order struct fields from largest alignment requirement to smallest (e.g. `int64` before `byte`) to minimize padding and improve data alignment.