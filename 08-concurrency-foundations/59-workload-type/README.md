# 08-concurrency-foundations / 59-workload-type

**Mistake #59: Ignoring workload type** — from *100 Go Mistakes and How to Avoid Them*, Chapter 8 (Concurrency Foundations).

## What this example shows

`main.go` compares two ways of processing bytes read from an `io.Reader`:

- `read1` — sequential: read a chunk, process it, repeat.
- `read2` — concurrent: a producer goroutine reads chunks and sends them over a buffered channel to a pool of 10 worker goroutines, which process chunks in parallel and aggregate the result with an atomic counter.

Both print `3072` when run against the `dummyReader` (which returns 3-byte chunks three times before `io.EOF`; note that `task` counts `len(b)` — the 1024-byte buffer size — rather than the actual bytes read, so `3 × 1024 = 3072`).

## The mistake

`read2` adds goroutines, a channel, a `WaitGroup`, and atomic operations to do work that is trivial (`len(b)`) and driven by a strictly sequential source (`io.Reader.Read` is not safe for concurrent use and yields one chunk at a time).

The lesson of mistake #59 is to **match your concurrency design to the workload type**:

| Factor | `read1` (sequential) | `read2` (concurrent pool) |
|--------|----------------------|---------------------------|
| Work per chunk | trivial | trivial |
| Source | sequential `io.Reader` | sequential `io.Reader` → fan-out |
| Synchronization | none | channel + `WaitGroup` + atomic |
| Best fit | I/O-bound, sequential, cheap-per-item | CPU-bound, expensive-per-item, parallelizable |

Concurrency is not free. When the per-item work is cheap and the input stream is inherently sequential, the coordination overhead of a worker pool can outweigh any parallelism benefit. Measure before parallelizing, and choose the model that fits the workload (I/O-bound vs CPU-bound, latency vs throughput, independent vs ordered items).

## Running

```bash
go run main.go
```

Expected output:

```
3072
3072
```

## Files

- `main.go` — sequential (`read1`) vs concurrent worker-pool (`read2`) implementations plus a `dummyReader` test source.
- `README.md` — this file.

## Further reading

- *100 Go Mistakes and How to Avoid Them*, Mistake 59: Ignoring workload type.
- Go blog: [Concurrency is not parallelism](https://go.dev/blog/concurrency-is-not-parallelism)
