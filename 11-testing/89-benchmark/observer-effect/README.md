# Observer Effect in Benchmarking

This example demonstrates **Mistake #89: Failing to account for the observer effect when benchmarking** from *100 Go Mistakes and How to Avoid Them*.

The observer effect occurs when the act of measuring changes the behavior being measured. In Go benchmarks, this commonly happens when the way you structure the benchmark (what you reuse, what you reallocate, what you time) changes the observed results — often in ways that have nothing to do with the code you intended to measure.

## What the code does

`main.go` defines two sum functions that operate on slices of fixed-size arrays:

- `calculateSum512(s [][512]int64)` — sums the first 8 elements of each row in a `[512]int64` matrix (4096 bytes per row, exactly one memory page).
- `calculateSum513(s [][513]int64)` — same logic, but with a `[513]int64` matrix (4104 bytes per row, spanning two pages).

`main_test.go` contains four benchmarks comparing the two array sizes:

| Benchmark | Matrix handling | Timing |
|-----------|-----------------|--------|
| `BenchmarkCalculateSum512_1` | Matrix created once, reused across all iterations | Sum only |
| `BenchmarkCalculateSum513_1` | Matrix created once, reused across all iterations | Sum only |
| `BenchmarkCalculateSum512_2` | Matrix created fresh each iteration (outside timer via `StopTimer`/`StartTimer`) | Sum only |
| `BenchmarkCalculateSum513_2` | Matrix created fresh each iteration (outside timer via `StopTimer`/`StartTimer`) | Sum only |

## The observer effect in action

The key difference between `_1` and `_2` is **cache behavior**, not just timing structure:

- In `_1`, the same matrix is reused for every iteration. After the first pass, the data stays hot in CPU cache, making the sum appear faster.
- In `_2`, a fresh matrix is allocated each iteration. Even though allocation happens outside the timed region (via `StopTimer`/`StartTimer`), each sum operates on cold, uncached memory, making it appear slower.

This is the observer effect: **the benchmark structure itself changes the observed performance.** Comparing `_1` against `_2` (or comparing 512 vs 513 across different structures) can lead you to conclude that one algorithm is faster than another, when the difference is really caused by cache locality introduced by how you wrote the benchmark.

There are also copy-paste bugs in the `_513` benchmarks — they call `createMatrix512` and `calculateSum512` instead of the `513` variants — which is another risk when you duplicate benchmark code without careful review.

## Running the benchmarks

```bash
go test -bench=. -benchmem
```

Use `-benchtime` to increase iteration counts and reduce noise:

```bash
go test -bench=. -benchmem -benchtime=1000x
```

Compare the `_1` vs `_2` results to see how much reuse-vs-reallocation (cache effects) distorts the timings.

## Lesson

- Use `b.StopTimer()` / `b.StartTimer()` to exclude setup work from the measured region.
- Be aware that reusing data across iterations warms the cache; reallocating each iteration leaves it cold. Either choice changes what you measure.
- Keep benchmarks consistent: don't mix structures when comparing variants.
- Review copied benchmarks carefully — it's easy to leave the wrong function name in a duplicate.
