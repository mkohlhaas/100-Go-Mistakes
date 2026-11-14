# 68 — String formatting while holding a lock

Part of the **100 Go Mistakes** concurrency practice series (see `09-concurrency-practice`).

## The mistake

`fmt.Errorf` (and `fmt.Sprintf`) with the `%v` verb on a value that implements `fmt.Stringer` calls that type's `String()` method. If `String()` acquires the same mutex that the caller already holds, the program deadlocks — Go's `sync.Mutex` / `sync.RWMutex` are **not reentrant**.

In this example `Customer` has:

- a `sync.RWMutex` guarding `id` and `age`
- a `String()` method that takes an `RLock`

`UpdateAge1` holds a write lock (`Lock`) and then builds an error message with `%v` on `c`, which invokes `String()` → `RLock` → deadlock.

## The code

```go
func (c *Customer) UpdateAge1(age int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if age < 0 {
		return fmt.Errorf("age should be positive for customer %v", c) // deadlock!
	}

	c.age = age
	return nil
}
```

Running `main()` with an invalid age (`-1`) deadlocks on the first call and never reaches `UpdateAge2` or `UpdateAge3`.

## Variants compared

| Method | Locking | Format string | Result |
|--------|---------|---------------|--------|
| `UpdateAge1` | locks before check | `%v` on `c` | **Deadlock** — `String()` tries to `RLock` while the write lock is held |
| `UpdateAge2` | validates before locking | `%v` on `c` | No deadlock, but check and write are not atomic (TOCTOU race) |
| `UpdateAge3` | locks before check | `%s` on `c.id` | **Correct** — no `String()` call, validation and write are atomic |

## The fix

Never format a `Stringer` with `%v` (or `%+v`, `%#v`) while holding its mutex. Either:

1. Validate before acquiring the lock (`UpdateAge2`), accepting that the check and mutation aren't atomic, or
2. Format only primitive fields with `%s`/`%d` while holding the lock (`UpdateAge3`), which is the correct approach when the check and mutation must be atomic.

## How to run

```bash
go run main.go
```

Expected behavior: the program deadlocks on `UpdateAge1` (you'll need to Ctrl+C it). Comment out the `UpdateAge1` call to see `UpdateAge2` and `UpdateAge3` run to completion.

## Key lesson

`%v` on a struct is not just formatting — it may invoke methods that touch shared state. When a mutex is held, any formatting that calls `String()`, `Error()`, or other methods on the guarded type can deadlock. Format primitive fields, or validate before locking.
