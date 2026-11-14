# Mutex-protected maps (Mistake #70)

Example from *100 Go Mistakes and How to Avoid Them*, chapter on concurrency — mistake **70: Mutex-protected slices and maps**.

## The code

`main.go` defines a `Cache` that stores account balances in a `map[string]float64` guarded by a `sync.RWMutex`, and computes the average balance three different ways:

| Method | What it does | Safe? |
|---|---|---|
| `AverageBalance1()` | Grabs a `RLock`, copies the **map reference**, releases the lock, then iterates the map | ❌ No — another goroutine may mutate the map while it is being iterated |
| `AverageBalance2()` | Holds the `RLock` for the entire iteration | ✅ Yes, but holds the read lock longer than necessary |
| `AverageBalance3()` | Copies the map entries while holding the `RLock`, releases the lock, then iterates the copy | ✅ Yes — recommended |

Run it:

```bash
go run main.go
```

## The mistake

Iterating over a map (or slice) after releasing the mutex that protects it is a data race: a concurrent writer can resize or modify the underlying map while the reader ranges over it, which can crash the program (`fatal error: concurrent map iteration and map write`) or silently produce wrong results.

Two safe options:

1. **Hold the read lock for the whole iteration** (`AverageBalance2`) — simple, but blocks writers for as long as the iteration takes.
2. **Copy the data under the lock, then iterate the copy without the lock** (`AverageBalance3`) — the recommended approach: readers never block each other, and writers are blocked only long enough to copy the entries.

## Related

- Parent folder: `09-concurrency-practice`
- Book: *100 Go Mistakes and How to Avoid Them* by Teiva Harsanyi (Mistake 70)
