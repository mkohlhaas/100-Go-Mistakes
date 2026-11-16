# Compiler Optimizations Benchmark

A small Go benchmark example measuring the effect of compiler optimizations on a
population count (popcount / Hamming weight) function.

## Files

- `main.go` — `popcnt(x uint64) uint64` using the SWAR (SIMD Within A Register)
  bit-counting algorithm with 64-bit masks.
- `main_test.go` — two benchmarks:
  - `BenchmarkPopcnt1` — calls `popcnt` but discards the result.
  - `BenchmarkPopcnt2` — stores the result in a package-level `global` variable,
    so the compiler cannot eliminate the call as dead code.

## Run the benchmarks

```bash
go test -bench=. -benchmem
```

Compare the two benchmarks to see how dead-code elimination (DCE) affects the
measured cost of the function.

## Inspect compiler optimizations

To see whether the compiler inlines or eliminates the call:

```bash
go test -gcflags="-m" -run=NONE -bench=BenchmarkPopcnt1
```

To view the generated assembly:

```bash
go test -c -gcflags="-l"   # build the test binary without inlining
go tool objdump main.test  # then inspect popcnt / the benchmark bodies
```

Or use `go tool compile -S` on `main.go` to see the machine code for `popcnt`
itself.

## Context

This is mistake #89 ("benchmarking compiler optimizations") from chapter 11
(Testing) of *100 Go Mistakes and How to Avoid Them*.