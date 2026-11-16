# Slice of Structs vs. Slice of Arrays

This example demonstrates a classic CPU-cache optimization: **only iterate over the data you actually need**, because unused fields pollute the cache and waste memory bandwidth.

## The problem

`Foo` is a struct with two `int64` fields, `a` and `b`. `sumFoo` sums only field `a`:

```go
type Foo struct {
    a int64
    b int64
}

func sumFoo(foos []Foo) int64 {
    var total int64
    for i := range foos {
        total += foos[i].a
    }
    return total
}
```

When the CPU loads `foos[i].a`, it pulls an entire **cache line** (typically 64 bytes) into L1. Since `Foo` is 16 bytes, a single cache line holds four elements — but each element carries 8 unused bytes of `b`. If `b` spans another cache line, you waste an entire line on data you never touch. Over a million elements that is dozens of megabytes of useless memory traffic.

## The fix

Split the fields into separate slices (SoA — structure of arrays). Field `b` exists but is **unused** by `sumBar`, so it lives in its own slice and is never touched:

```go
type Bar struct {
    a []int64
    b []int64
}

func sumBar(bar Bar) int64 {
    var total int64
    for i := range bar.a {
        total += bar.a[i]
    }
    return total
}
```

Now `bar.a` is a tight, contiguous `[]int64`. Every loaded cache line contains only values that are actually summed — the unused `b` slice never pollutes the cache — so the cache holds more useful data and fewer lines are fetched.

## Run the benchmarks

```bash
go test -bench=. -benchmem
```

Expected outcome: `BenchmarkSumBar` should be faster and allocate no more than `BenchmarkSumFoo` (both allocate only for setup), because it touches less memory per element.

## Reference

This is mistake **91** from *[100 Go Mistakes and How to Avoid Them](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them)* — chapter "12 - Optimizations", section on CPU caches.