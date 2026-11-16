# 93 — Instruction-level parallelism

Example from *[100 Go Mistakes and How to Avoid Them]*, optimization #93:
prefer local variables over repeated memory/array access inside hot loops so
the CPU can execute independent iterations in parallel (instruction-level
parallelism).

## The code

`main.go` contains two functionally equivalent functions:

```go
func add(s [2]int64) [2]int64 {
	for range n {
		s[0]++
		if s[0]%2 == 0 {
			s[1]++
		}
	}
	return s
}

func add2(s [2]int64) [2]int64 {
	for range n {
		v := s[0]
		s[0] = v + 1
		if v%2 != 0 {
			s[1]++
		}
	}
	return s
}
```

## Why `add2` is faster

In `add`, every iteration reads and writes `s[0]` from the array. The next
iteration cannot start until the current one finishes writing `s[0]`, because
the condition `s[0]%2 == 0` depends on that written value. This creates a
loop-carried dependency that serializes execution and limits instruction-level
parallelism.

In `add2`, the value is loaded once into a local variable `v`. The increment
and the parity check operate on `v`, which lives in a register, so the CPU can
overlap multiple loop iterations — it does not have to wait for `s[0]` to be
written back to memory before starting the next iteration.

## Running the benchmarks

```sh
go test -bench=. -benchmem
```

### Results (representative)

Results vary by machine; the two functions are essentially equivalent on this
hardware (Intel Core i5-5250U, amd64), each taking roughly 2.3 ms per iteration
with zero allocations:

```
goos: linux
goarch: amd64
pkg: 100-go-mistakes/12-optimizations/93-instruction-level-parallelism
cpu: Intel(R) Core(TM) i5-5250U CPU @ 1.60GHz
BenchmarkAdd-4     516    2308881 ns/op    0 B/op    0 allocs/op
BenchmarkAdd2-4    518    2320761 ns/op    0 B/op    0 allocs/op
PASS
ok  	100-go-mistakes/12-optimizations/93-instruction-level-parallelism	2.399s
```

Neither version meaningfully outperforms the other here, so this example mainly
serves to illustrate the *pattern* (load once into a local, then work on the
register) rather than guaranteeing a speedup. Actual gains depend on the CPU,
compiler version, and how tight the dependency chain is.

[100 Go Mistakes and How to Avoid Them]: https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them