# timer

Example comparing two ways to handle setup in Go benchmarks.

## The code

`main_test.go` contains two benchmarks for the same work:

```go
func BenchmarkFoo1(b *testing.B) {
	expensiveSetup()
	b.ResetTimer()
	for b.Loop() {
		functionUnderTest()
	}
}

func BenchmarkFoo2(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		expensiveSetup()
		b.StartTimer()
		functionUnderTest()
	}
}
```

## The difference

| Benchmark | Setup timing | Meaning |
|---|---|---|
| `BenchmarkFoo1` | Runs **once**, before the loop starts | `ResetTimer()` discards elapsed time after one-time setup, so only `functionUnderTest()` is measured. This is correct when setup is not part of the operation being benchmarked. |
| `BenchmarkFoo2` | Runs **every iteration** | `StopTimer()`/`StartTimer()` inside the loop pauses timing around `expensiveSetup()` on every iteration. The setup cost is excluded from measurement but still runs `b.N` times, wasting CPU. If setup truly belongs to each iteration, measure it inside the timer instead. |

## Run it

```bash
go test -bench=. -benchmem
```

The default run benchmarks each function for about 1 second of CPU time, so it can take several seconds overall. For a single quick iteration per benchmark:

```bash
go test -bench=. -benchtime=1x -timeout=10s
```

Example output (with `-benchtime=1x`):

```text
goos: linux
goarch: amd64
pkg: 100-go-mistakes/11-testing/89-benchmark/timer
cpu: Intel(R) Core(TM) i5-5250U CPU @ 1.60GHz
BenchmarkFoo1-4   	       1	       359.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkFoo2-4   	       1	       668.0 ns/op	       0 B/op	       0 allocs/op
PASS
```

## Takeaway

Use `b.ResetTimer()` after one-time setup performed before the loop. Avoid putting setup inside the loop with `StopTimer()`/`StartTimer()` unless the setup is genuinely part of each operation being measured — otherwise you pay its cost `b.N` times without measuring it.

Reference: Chapter 11, Mistake 89 — "Benchmarking" from [100 Go Mistakes](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them).
