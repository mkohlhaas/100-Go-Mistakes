# Predictability Demo

Linked lists are a performance killer in any programming language. This example shows why: they force the CPU to **chase pointers**, which defeats hardware prefetching and branch prediction, while a flat slice with a predictable access pattern lets the CPU run at full speed.

## The Code

`main.go` contains two functions that sum `int64` values:

- `linkedList` — walks a singly-linked list of `node` structs, following `next` pointers until `nil`
- `sum2` — reads every 2nd element of an `int64` slice (`i += 2`)

`main_test.go` benchmarks both over 1,000,000 elements.

## Why Linked Lists Are Slow

CPUs are fastest when memory access is **predictable**.

`sum2` walks the slice with a fixed stride and a known trip count. The hardware prefetcher sees the pattern coming and pulls the next cache lines into L1 before the loop needs them. The loop is also friendly to branch prediction since the iteration count is fixed.

`linkedList` is **pointer chasing**: each iteration loads a node, then uses the address stored inside it (`n.next`) to find the next one. The CPU cannot know where the next node lives until it has loaded the current one, so the prefetcher is blind and the loop stalls waiting on memory latency. Every hop is a serial dependency on main memory.

That is why traversing a hand-linked list is typically much slower than traversing a flat slice — even when both do the same amount of arithmetic.

## Run the Benchmarks

```bash
go test -bench=. -benchmem .
```

Example output:

```
BenchmarkLinkedList-8    10000000    150 ns/op    0 B/op    0 allocs/op
BenchmarkSum2-8          20000000     80 ns/op    0 B/op    0 allocs/op
```

Exact numbers vary by CPU, but `sum2` is typically faster because its access pattern is predictable and prefetchable, while the linked list is not.

## The Caveat — Real Life Is Even Worse

This benchmark allocates the nodes as one contiguous slice (`make([]node, n)`) and links them in order, so the nodes happen to sit next to each other in memory. That understates the problem: in real Go code, linked-list nodes are usually individual heap allocations scattered across the heap. Each hop then jumps to an unrelated address, making pointer chasing even slower than this test already shows.

Also, `sum2` only sums half the elements, so it does less work than a full traversal. The point is not that stride-2 is magic — it's that **predictable, linear access over a dense slice lets the CPU do its job**, while data-dependent pointer chasing defeats prefetching and branch prediction.

## Bottom Line

Prefer slices over linked lists whenever you only need to iterate. If you're building a graph or tree where pointers are unavoidable, keep hot traversal data contiguous (e.g., store nodes in a slice and refer to them by index) so the CPU can prefetch.

## Context

Part of **Chapter 12 — Optimizations**, Mistake #91: "Ignoring CPU cache predictability" in *100 Go Mistakes and How to Avoid Them*.
