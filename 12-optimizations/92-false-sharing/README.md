# False Sharing

Example from *100 Go Mistakes and How to Avoid Them*, mistake #92:
**ignoring false sharing**.

## The problem

`Result1` holds two `int64` counters:

```go
type Result1 struct {
	sumA int64
	sumB int64
}
```

Two goroutines each increment a different field concurrently. Because both
fields live on the same 64-byte CPU cache line (16 bytes total), the cores
constantly invalidate each other's copy of the line as they write. This
cache-line ping-pong is **false sharing**: the fields are logically
independent, but the hardware treats them as contended.

## The fix

`Result2` inserts padding between the fields so they land on separate cache
lines:

```go
type Result2 struct {
	sumA int64
	_    [56]byte
	sumB int64
}
```

`sumA` occupies bytes 0–7, the padding spans bytes 8–63, and `sumB` starts at
byte 64 — the next cache line. The two goroutines no longer contend for the
same line.

## Run the benchmarks

```bash
go test -bench=. -benchmem
```

Expected result: `BenchmarkCount2` runs faster and allocates no more than
`BenchmarkCount1`; the difference is most visible on multi-core machines.

Example output (Intel Core i5-5250U, Go 1.27.4):

```text
goos: linux
goarch: amd64
pkg: 100-go-mistakes/12-optimizations/92-false-sharing
cpu: Intel(R) Core(TM) i5-5250U CPU @ 1.60GHz
BenchmarkCount1-4   	     321	   3736783 ns/op	     270 B/op	       4 allocs/op
BenchmarkCount2-4   	     415	   3405894 ns/op	     230 B/op	       4 allocs/op
```

Your exact numbers will vary by CPU and core count. The effect is most visible on
machines with many cores and separate per-core caches under heavy contention; on
a dual-core or lightly-contended machine the gap can be small.

## Further reading

- "100 Go Mistakes and How to Avoid Them", Teiva Harsanyi, mistake #92
- [Go memory model](https://go.dev/ref/mem)
