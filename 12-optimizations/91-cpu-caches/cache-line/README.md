# Cache Line Demo

This example demonstrates how memory access patterns interact with CPU **cache lines** — a classic Go performance mistake.

## The Code

`main.go` contains two functions that sum elements of an `int64` slice at different strides:

- `sum2` — reads every 2nd element (`i += 2`)
- `sum8` — reads every 8th element (`i += 8`)

`main_test.go` benchmarks both over a 1,000,000-element slice.

## Why It Matters

CPUs load memory in **cache lines** — typically 64 bytes on x86. Since `int64` is 8 bytes, one cache line holds 8 `int64` values.

When you read one element from a slice, the CPU brings in the whole cache line containing it. If your loop stride skips most of those values (like `i += 8`), you waste memory bandwidth fetching data you never use. Accessing elements contiguously (or with a small stride) makes better use of each cache line and is usually faster.

This is why iterating a slice sequentially is far quicker than jumping around — even when you do the same amount of arithmetic.

## Run the Benchmarks

```bash
go test -bench=. -benchmem .
```

Example output:

```
BenchmarkSum2-8     10000000    123 ns/op    0 B/op    0 allocs/op
BenchmarkSum8-8      5000000    234 ns/op    0 B/op    0 allocs/op
```

Exact numbers vary by CPU, but `sum2` typically beats `sum8` because it reuses more of each cache line.

## The learning point (and a caveat)

This example is intentionally simple, but it has a flaw: `sum8` reads 8× fewer elements than `sum2`, so any speed difference isn't purely a cache-line effect. If your algorithm genuinely needs every 8th element, you don't have a choice about the stride — you either accept the memory-bandwidth cost or change the problem/data layout.

The real lessons of "ignoring CPU cache lines" are about **data layout**, not about picking a stride:

### 1. Store hot fields separately (SoA vs AoS)

If you're striding because you're filtering a slice of structs:

```go
type Item struct {
    Active bool
    Value  int64
    Meta   [6]int64 // cold data
}
items []Item
for i := 0; i < len(items); i += 8 {
    total += items[i].Value
}
```

you load 64-byte cache lines full of `Meta` and `Active` just to read one `Value`. The fix is to store the hot field separately:

```go
values []int64 // dense, contiguous
for i := 0; i < len(values); i++ {
    total += values[i]
}
```

Now every cache line you fetch contains 8 values you actually use.

### 2. Avoid false sharing

If multiple goroutines each write to their own `int64` counter, but those counters live on the same cache line, the CPU cores ping-pong that line back and forth even though they're logically independent. Give each hot field its own cache line:

```go
type Counter struct {
    value int64
    _     [56]byte // pad to 64 bytes
}
```

### 3. Block / tile your iteration

When walking a 2D matrix or nested structure, process cache-line-sized blocks (e.g., 8×8 tiles for `int64`) so you reuse fetched lines before they're evicted, rather than striding across rows.

### Bottom line

You can't "choose" to stride by 2 if the algorithm needs stride 8. The choice is upstream — how you lay out data in memory, and whether you're paying to load bytes you never touch.

## Context

Part of **Chapter 12 — Optimizations**, Mistake #91: "Ignoring CPU cache lines" in *100 Go Mistakes and How to Avoid Them*.
