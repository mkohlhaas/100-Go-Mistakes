# Slice Memory Leak

This example demonstrates **Mistake #26** from *100 Go Mistakes*: reslicing a Go slice does not copy data or shrink the underlying allocation. It only creates a new slice header pointing at the same backing array. If you hold onto a small sub-slice of a large slice, you unintentionally keep the entire backing array — and everything it references — alive in memory.

## What the code does

The program allocates a slice of 1,000 `Foo` structs, each holding a 1 MB `[]byte` (~1 GB total). It then keeps only the first two elements using three different approaches and reports memory usage before and after `runtime.GC()`.

## Scenarios

### Scenario 1 — `keepFirstTwoElementsOnly` (reslice) — **leaks memory**

```go
func keepFirstTwoElementsOnly(foos []Foo) []Foo {
    return foos[:2]
}
```

`foos[:2]` shares the original backing array. The garbage collector cannot free any of the 1,000 byte slices because the returned slice keeps the whole array reachable. Even after `runtime.GC()`, ~1 GB stays pinned.

### Scenario 2 — `keepFirstTwoElementsOnlyCopy` (copy) — **fixes it**

```go
func keepFirstTwoElementsOnlyCopy(foos []Foo) []Foo {
    res := make([]Foo, 2)
    copy(res, foos)
    return res
}
```

Allocates a new slice and copies only the first two elements into it. The original backing array becomes unreachable and is freed by the GC. After GC, only ~2 KB remains.

### Scenario 3 — `keepFirstTwoElementsOnlyMarkNil` (mark nil) — **fixes it**

```go
func keepFirstTwoElementsOnlyMarkNil(foos []Foo) []Foo {
    for i := 2; i < len(foos); i++ {
        foos[i].v = nil
    }
    return foos[:2]
}
```

Explicitly nils out the `v` fields of elements 2–999 before reslicing. This breaks the references so the unused byte slices can be collected. After GC, the unused elements are freed.

## The takeaway

Reslicing never copies data. To truly release memory from a large slice when you only need a small portion, either:

1. **Copy** the needed elements into a new slice, or
2. **Nil out** the references you no longer need before reslicing.

## Note on liveness analysis

Earlier versions of this example explicitly nil-ed the variables at the end of each scenario:

```go
two = nil
foos = nil
runtime.GC()
```

These assignments were removed because they are **ineffectual** — Go's compiler performs liveness analysis and treats `two` and `foos` as dead once they are no longer read (after `runtime.KeepAlive(two)`). The backing array is released by the GC based on reachability at that point, regardless of whether you explicitly nil the local variables. `staticcheck` (SA4006/SA4016) flags such assignments as dead code, so they were dropped.

The only thing that genuinely matters for releasing memory is **not holding onto references** you no longer need — either by copying out what you want (Scenario 2), by nil-ing out the fields you don't want before reslicing (Scenario 3), or simply by letting a reslice variable like `two` fall out of scope (Scenario 1).

### Why `runtime.KeepAlive(two)` is present

Each scenario calls `runtime.KeepAlive(two)` because `two` is otherwise never read — without it, Go reports `declared and not used` and the program won't compile:

```go
two := keepFirstTwoElementsOnly(foos)
// ... nothing reads 'two' ...
```

`runtime.KeepAlive(two)` is the sole read of `two`, and it anchors `two`'s liveness through the first `runtime.GC()`. That matters most in Scenario 1: without `KeepAlive`, the compiler could treat `two` as dead immediately after `two := keepFirstTwoElementsOnly(foos)` (since it's never read), the GC might collect the 1 GB backing array before we observe it, and we'd never see the leak.

## Run it

```bash
go run main.go
```

You will see the allocated memory stay high in Scenario 1 after GC, while Scenarios 2 and 3 release it.
