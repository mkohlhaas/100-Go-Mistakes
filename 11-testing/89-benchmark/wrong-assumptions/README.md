# wrong-assumptions

A Go microbenchmark demonstrating a wrong assumption about benchmarking: **the compiler can eliminate stores whose results are never observed**.

## What it does

`main_test.go` contains two benchmarks:

- `BenchmarkAtomicStoreInt32` — stores `1` into an `atomic.Int32`
- `BenchmarkAtomicStoreInt64` — stores `1` into an `atomic.Int64`

Both write to a local variable that is never read, so the stores are dead: nothing observes the value afterward. The compiler (and the testing harness's optimization passes) is free to remove them entirely.

## The wrong assumption

A common mistake is assuming that benchmarking an operation measures the cost of that operation:

```go
var v atomic.Int32
for b.Loop() {
    v.Store(1) // WRONG: v is never read; this may be optimized away
}
```

Because `v` escapes nothing and its value is never used, the loop body can be dead-code eliminated. The benchmark may then measure nothing (or near-nothing), giving misleadingly fast results that have nothing to do with the actual cost of `atomic.Int32.Store` vs. `atomic.Int64.Store`.

## How to run

```sh
go test -bench=. -benchmem
```

Compare the results before and after making the stored value observable (for example, by summing `v` at the end or assigning it to a package-level sink) to see how the numbers change when the compiler can no longer eliminate the stores.

## Lesson

Always make sure the work you benchmark is actually performed — use the result, return it, or sink it — otherwise you're benchmarking the optimizer, not your code.
