# Stack vs Heap

Example from *100 Go Mistakes and How to Avoid Them*, mistake #95:
**returning a pointer to a local variable**.

## The problem

Go allocates variables either on the stack (fast, reclaimed when the function
returns) or on the heap (slower, garbage collected). The compiler's escape
analysis decides which one: if a variable's address escapes the function — for
example because a pointer to it is returned — it must be heap-allocated.

In `listing2`, `sumPtr` returns `&z`, where `z` is a local variable. Since the
pointer outlives the stack frame, `z` escapes to the heap:

```go
//go:noinline
func sumPtr(x, y int) *int {
	z := x + y
	return &z // z escapes to the heap
}
```

Every call to `sumPtr` allocates 8 bytes on the heap, which the garbage
collector must later trace and free.

## The fix

Return values by value when you can. `sumValue` returns an `int` by value, so
nothing escapes and the compiler keeps everything on the stack:

```go
//go:noinline
func sumValue(x, y int) int {
	z := x + y
	return z // no escape: stays on the stack
}
```

If a caller needs the result's address, have it take the address of the returned
value instead:

```go
c := sumValue(a, b)
ptr := &c // c lives in the caller's stack frame; &c only escapes if it must
```

Whether `c` itself lands on the stack or heap then depends on how the caller
uses it — not on `sumValue`.

`listing3`'s `sum` takes pointer arguments but returns an `int` by value, so no
heap allocation is needed inside the function either.

## Run the benchmarks

```bash
go test -bench=. -benchmem
```

`BenchmarkSumValue` reports `0 B/op` and `0 allocs/op`, while
`BenchmarkSumPtr` reports `8 B/op` and `1 allocs/op` because `z` escapes to the
heap.

Example output (Intel Core i5-5250U, Go 1.27.4):

```text
goos: linux
goarch: amd64
pkg: 100-go-mistakes/12-optimizations/95-stack-heap
cpu: Intel(R) Core(TM) i5-5250U CPU @ 1.60GHz
BenchmarkSumValue-4   	428865654	         2.770 ns/op	       0 B/op	       0 allocs/op
BenchmarkSumPtr-4     	63854835	        38.25 ns/op	       8 B/op	       1 allocs/op
PASS
ok  	100-go-mistakes/12-optimizations/95-stack-heap	3.635s
```

Your exact numbers will vary. The deterministic signal is the allocation count:
returning a pointer to a local costs one heap allocation per call; returning a
value costs none.

### Why the heap version is slower

| Benchmark | ns/op | B/op | allocs/op |
|---|---:|---:|---:|
| `SumValue` — returns `int` by value (stack) | ~2.7 | 0 | 0 |
| `SumPtr` — returns `*int` to a local (heap) | ~35 | 8 | 1 |

`SumPtr` is roughly **13× slower** in this run. Every call pays for a heap
allocation (8 bytes) plus the garbage-collection work to trace and free it
later. The stack version keeps `z` in the caller's frame with no allocation at
all.

## Rule of Thumb

Prefer returning values by value. Only return a pointer when the caller genuinely
needs shared or mutable state — and remember that doing so forces the pointee
onto the heap.

## Further reading

- "100 Go Mistakes and How to Avoid Them", Teiva Harsanyi, mistake #95
- [Go escape analysis](https://go.dev/doc/go1.9#escape_analysis)
- [`go build -gcflags="-m"` to see escape decisions](https://go.dev/doc/gcguide)
