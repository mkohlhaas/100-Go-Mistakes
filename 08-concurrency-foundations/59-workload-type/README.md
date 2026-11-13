# 08-concurrency-foundations / 59-workload-type

**Mistake #59: Ignoring workload type** — from *100 Go Mistakes and How to Avoid Them*, Chapter 8 (Concurrency Foundations).

## What this example shows

`main.go` compares two ways of processing chunks read from an `io.Reader`:

- `read1` — sequential: read a chunk, process it, repeat.
- `read2` — concurrent: a single producer goroutine reads chunks and sends them over a buffered channel to a pool of 10 worker goroutines, which process chunks in parallel and aggregate the result with an atomic counter.

Both print `5068800` when run against the `dummyReader` (which returns 100 chunks of 1024 bytes each).

## The mistake

`read2` adds goroutines, a channel, a `WaitGroup`, and atomic operations to do work that could be done sequentially. Whether that trade-off is worth it depends entirely on the **workload type**:

| Factor | `read1` (sequential) | `read2` (worker pool) |
|--------|----------------------|---------------------------|
| Work per chunk | trivial / cheap | CPU-bound, expensive |
| Source | sequential `io.Reader` | sequential `io.Reader` → fan-out |
| Concurrency safe? | Yes (single goroutine) | Yes (only producer calls `Read`) |
| Ordering guaranteed | Yes | No (results aggregate out of order) |
| Synchronization | none | channel + `WaitGroup` + atomic |
| Best fit | I/O-bound, cheap-per-item, ordered | CPU-bound, independent, parallelizable |

`io.Reader` is not safe for concurrent use, so only one goroutine may call `Read` — which is why `read2` uses a producer-consumer pattern rather than multiple readers.

## Key fixes over the original snippet

- **`b[:n]` instead of `b`** — both functions process only the bytes actually returned by `Read`, not the entire 1024-byte buffer.
- **`errors.Is(err, io.EOF)`** — correctly detects EOF even if a reader wraps the error (prefer this over `err == io.EOF`).
- **CPU-bound `task`** — `task` now does real per-byte computation so that parallelism across workers can actually reduce runtime, given enough chunks to amortize goroutine overhead.

## When concurrency helps (and when it doesn't)

Concurrency is not free. A worker pool only pays off when:

1. The per-chunk work is CPU-bound enough to offset goroutine/channel overhead.
2. There are enough chunks to keep all workers busy.
3. Chunks are independent (no required ordering).

If `task()` is trivial (e.g. `len(b)`), `read1` is almost always the better choice: less code, no synchronization bugs, and no coordination tax.

## Running

```bash
go run main.go
```

Expected output:

```
5068800
5068800
```

## Files

- `main.go` — sequential (`read1`) vs concurrent worker-pool (`read2`) implementations plus a `dummyReader` test source.
- `README.md` — this file.

## Further reading

- *100 Go Mistakes and How to Avoid Them*, Mistake 59: Ignoring workload type.
- Go blog: [Concurrency is not parallelism](https://go.dev/blog/concurrency-is-not-parallelism)
