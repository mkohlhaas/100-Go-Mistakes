# Copying `sync` types (Mistake #74)

Example from *100 Go Mistakes and How to Avoid Them*, chapter on concurrency — mistake **74: Copying a `sync.Mutex` (or any `sync` type) by value**.

## The code

`main.go` defines two counters that guard a `map[string]int` with a `sync.Mutex`:

| Type | Mutex field | Receiver | Safe? |
|---|---|---|---|
| `Counter` | `mu sync.Mutex` (by value) | value receiver `(c Counter)` | ❌ No |
| `Counter2` | `mu *sync.Mutex` (pointer) | pointer receiver `(c *Counter2)` | ✅ Yes |

`main` launches two goroutines that each call `Increment1/2` on a `Counter1/2` value:

```go
counter := NewCounter()

go func() {
    counter.Increment1("foo")
}()
go func() {
    counter.Increment1("bar")
}()

time.Sleep(10 * time.Millisecond)
```

## The mistake

`sync.Mutex` (and `sync.RWMutex`, `sync.WaitGroup`, `sync.Once`, `sync.Pool`, `sync.Map`, `sync.Cond`, `atomic.*` types, channels, etc.) **must not be copied**. A mutex contains internal state (an internal semaphore/waiter queue); copying the struct copies that state, so every copy has its own, independent lock.

Here the damage is compounded by the **value receiver**:

```go
func (c Counter) Increment1(name string) {
    c.mu.Lock()          // locks the copy's mutex, not the original's
    defer c.mu.Unlock()
    c.counters[name]++   // map is a reference type: both copies write to the same map
}
```

When `counter.Increment1("foo")` is called, `c` is a copy of `Counter`. `c.mu` is a copy of the mutex, so `Lock`/`Unlock` operate on a mutex that no other goroutine ever touches. The map itself is a reference type, so both copies write to the same underlying map — but without mutual exclusion. The two goroutines therefore increment `counters["foo"]` / `counters["bar"]` concurrently — a **data race** on the shared map.

Run with the race detector to see it:

```bash
go run -race main.go
```

Typical output:

```
WARNING: DATA RACE
Write at 0x... by goroutine 6:
  main.(*Counter).Increment1(...)
    ...
Previous write at 0x... by goroutine 7:
  ...
```

## The fix

Never copy a value that contains a `sync` type. Either:

1. **Store a pointer to the mutex** (`mu *sync.Mutex`) — copying the struct copies the pointer, which still refers to the same mutex, so all copies synchronize on the same lock. This is what `Counter2` does.

   ```go
   type Counter2 struct {
       mu       *sync.Mutex // good
       counters map[string]int
   }

   func NewCounter2() Counter2 {
       return Counter2{
           mu:       &sync.Mutex{},
           counters: map[string]int{},
       }
   }

   func (c *Counter2) Increment2(name string) {
       c.mu.Lock()
       defer c.mu.Unlock()
       c.counters[name]++
   }
   ```

2. **Pass the struct by pointer** everywhere (`*Counter`) and use pointer receivers, so the value is never copied at all.

Go tooling enforces this: `go vet` runs the `copylocks` checker and reports something like:

```
main.go:30:9: Increment1 passes lock by value: ...Counter contains sync.Mutex
```

## Key rules

1. **Never copy `sync` types** — `sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, `sync.Once`, `sync.Pool`, `sync.Map`, `sync.Cond`, channels, `atomic` types, etc.
2. **Prefer `*sync.Mutex`** if the struct itself must be copied or passed by value.
3. **Otherwise pass the struct by pointer** and use pointer receivers so nothing is copied.
4. **Run `go vet`** (it includes `copylocks`) and `go run -race` to catch these mistakes.

## Running

```bash
go run main.go
go run -race main.go   # demonstrates the data race on the buggy Counter
go vet ./...
```

## Related

- Parent folder: `09-concurrency-practice`
- Book: *100 Go Mistakes and How to Avoid Them* by Teiva Harsanyi (Mistake 74)
- Previous topic in this chapter: `73-errgroup`, `72-cond`, `70-mutex-slices-maps`
