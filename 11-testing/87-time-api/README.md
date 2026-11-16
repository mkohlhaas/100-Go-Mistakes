# 87-time-api — Testing Code That Relies on `time.Now()`

This example comes from item **87** ("Beware of the time API") of the testing chapter in *100 Go Mistakes and How to Avoid Them*.

It shows how a `Cache` that trims old events evolves across four listings, each addressing the same core problem:

> **Functions that call `time.Now()` internally are hard to test deterministically** — the test outcome depends on wall-clock time, making tests flaky and timing-sensitive.

---

## Project structure

```
87-time-api/
├── listing1/   # The mistake: time.Now() called inside the method
├── listing2/   # First fix attempt: inject a now() function field (incomplete)
├── listing3/   # Second fix: pass now time.Time as a parameter
├── listing4/   # Best fix: pass the already-computed threshold time
└── README.md   # this file
```

Each listing contains `main.go` (implementation) and `main_test.go` (tests).

---

## The problem — listing 1

`listing1/main.go` calls `time.Now()` directly inside `TrimOlderThan` — flagged with `// ⚠️ calling time.Now()`:

```go
func (c *Cache) TrimOlderThan(since time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // ⚠️ calling time.Now()
    t := time.Now().Add(-since)
    for i := 0; i < len(c.events); i++ {
        if c.events[i].Timestamp.After(t) {
            c.events = c.events[i:]
            return
        }
    }
}
```

The test in `listing1/main_test.go` builds events relative to `time.Now()` — flagged with `// ⚠️ calling time.Now()`:

```go
events := []Event{
    // ⚠️ calling time.Now()
    {Timestamp: time.Now().Add(-20 * time.Millisecond)},
    {Timestamp: time.Now().Add(-10 * time.Millisecond)},
    {Timestamp: time.Now().Add(10 * time.Millisecond)},
}
```

**Why this is bad:**
- The test is timing-dependent. Scheduling delays, CI load, or a slow machine can flip the result.
- It is impossible to test edge cases (exactly at the boundary, far in the past/future) without sleeping or manipulating the system clock.

Run the tests and observe the flakiness:

```bash
go test ./listing1/
```

---

## Fix attempt #1 — listing 2

`listing2/main.go` adds a `now func() time.Time` field to `Cache` so tests can inject a fake clock:

```go
type now func() time.Time

type Cache struct {
    mu     sync.RWMutex
    events []Event
    now    now
}

func NewCache() *Cache {
    return &Cache{
        events: make([]Event, 0),
        now:    time.Now,
    }
}
```

The test injects a fixed time:

```go
cache := &Cache{now: func() time.Time {
    return parseTime(t, "2020-01-01T12:00:00.06Z")
}}
```

The catch: `TrimOlderThan` still calls `time.Now()` directly — flagged with `// ⚠️ still calling time.Now()` — so the injected `now` field is unused:

```go
func (c *Cache) TrimOlderThan(since time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // ⚠️ still calling time.Now()
    t := time.Now().Add(-since)
    // ...
}
```

---

## Fix #2 — listing 3: inject `now time.Time` as a parameter

Instead of generating time inside the method, the caller passes it in:

```go
func (c *Cache) TrimOlderThan(now time.Time, since time.Duration) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    // ⚠️ cache does time manipulation (not its business)
    t := now.Add(-since)
    // ...
}
```

The test now fully controls time — also flagged with `// ⚠️ cache does time manipulation (not its business)`:

```go
// ⚠️ cache does time manipulation (not its business)
cache.TrimOlderThan(parseTime(t, "2020-01-01T12:00:00.06Z"), 15*time.Millisecond)
```

**Pros:** deterministic, no sleeps, no wall clock.  
**Cons:** the cache method is still responsible for computing `now.Add(-since)`, mixing clock access with business logic.

---

## Best fix — listing 4: pass the computed threshold time

Move time computation entirely out of the method. `TrimOlderThan` receives the cutoff `t` directly:

```go
func (c *Cache) TrimOlderThan(t time.Time) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    for i := 0; i < len(c.events); i++ {
        if c.events[i].Timestamp.After(t) {
            c.events = c.events[i:]
            return
        }
    }
}
```

The caller decides what "now" means:

```go
cache.TrimOlderThan(parseTime(t, "2020-01-01T12:00:00.06Z").
    Add(-15 * time.Millisecond))
```

**Pros:**
- Fully deterministic tests.
- The cache knows nothing about clocks — it only compares timestamps.
- Most flexible: the caller can use a real clock, a mock clock, or any computed time.

---

## How to run

Each listing is its own Go package:

```bash
go test ./listing1/
go test ./listing2/
go test ./listing3/
go test ./listing4/
```

To run all at once from the `87-time-api` directory:

```bash
go test ./...
```

Note: because every listing uses `package listing1`, you cannot compile them into a single binary, but `go test ./...` will test each directory independently.

---

## Key takeaway

Never call `time.Now()` (or `time.Now().Add(...)`) inside code you need to test. Pass time in — either as a `time.Time` parameter or an injected clock function — and ideally pass the **already-computed** time so your business logic stays free of clock concerns. Listing 4 is the recommended approach.
