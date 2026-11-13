# 54 — Handle errors from deferred calls

Example for mistake **#54** from *100 Go Mistakes and How to Avoid Them*: ignoring errors returned by deferred functions (here, `*sql.Rows.Close`).

## The problem

```go
rows, err := db.Query(query, clientID)
if err != nil {
    return 0, err
}
defer rows.Close()   // error from Close() is silently discarded
```

`rows.Close()` can fail (e.g. a network or cursor error while releasing server-side resources).  
`defer rows.Close()` discards that error, so failures go unnoticed and resources may leak or be left in an inconsistent state.

## What's in `main.go`

| Function | Approach | Drawback |
|----------|----------|----------|
| `getBalance1` | `defer rows.Close()` | Error from `Close()` is silently ignored. |
| `getBalance2` | `defer func() { _ = rows.Close() }()` | Same as above, just explicit about ignoring. |
| `getBalance3` | Named returns + deferred closure | Captures both the original `err` and `Close()`'s error, logging or propagating them appropriately. |

## Recommended pattern

Use named return values so the deferred closure can inspect and combine both errors:

```go
func getBalance(db *sql.DB, clientID string) (balance float32, err error) {
    rows, err := db.Query(query, clientID)
    if err != nil {
        return 0, err
    }
    defer func() {
        if closeErr := rows.Close(); closeErr != nil {
            if err == nil {
                err = closeErr   // surface the close error if nothing else failed
            } else {
                log.Printf("failed to close rows: %v", closeErr)
            }
        }
    }()

    // use rows ...
    return 0, nil
}
```

## What's idiomatic?

There isn't one single "most idiomatic" answer — it depends on whether the close error is equally important to the caller or just diagnostic noise. Here are the three common patterns and when to use each.

### 1. Log the close error, return the primary error (very common for deferred cleanup)

```go
defer func() {
    if closeErr := rows.Close(); closeErr != nil {
        if err == nil {
            err = closeErr   // nothing else failed, so surface the close error
        } else {
            log.Printf("failed to close rows: %v", closeErr)
        }
    }
}()
```

This is arguably the most common idiom for `defer` cleanup when an earlier operation already failed. The caller gets the original, actionable error (e.g. "query failed"), and the close failure is recorded in logs for operators without cluttering the error contract.

### 2. Return only the close error if there's no primary error

If `db.Query` succeeded but `rows.Close()` fails, you should return `closeErr` so the caller knows cleanup didn't happen:

```go
defer func() {
    if closeErr := rows.Close(); closeErr != nil && err == nil {
        err = closeErr
    }
}()
```

### 3. Combine both errors (when both matter to the caller)

If the caller genuinely needs to inspect both failures, use `errors.Join` (Go 1.20+):

```go
defer func() {
    if closeErr := rows.Close(); closeErr != nil {
        err = errors.Join(err, closeErr)
    }
}()
```

Or with `fmt.Errorf` and `%w` on older Go versions, though `errors.Join` is cleaner for multiple peers.

## Rule of thumb

| Situation | Idiomatic choice |
|-----------|-----------------|
| Primary operation already failed | Log close error, return original error |
| Primary succeeded, cleanup failed | Return close error |
| Both errors are independently actionable | `errors.Join(err, closeErr)` |

So the `log.Printf` version in `getBalance3` is idiomatic because the query error is usually what the caller cares about; the close error is secondary diagnostics. Wrapping would be appropriate if callers need to distinguish or retry based on the close failure specifically.

## Possible ways to wrap an error

| Approach | Preserves chain? | Example | Use case |
|----------|-----------------|---------|----------|
| `fmt.Errorf(..., %w)` | Yes | `fmt.Errorf("read %s: %w", path, err)` | Add context to one error; callers can use `errors.Is` / `errors.As` |
| `errors.Join(...)` (Go 1.20+) | Yes | `errors.Join(err1, err2)` | Combine multiple independent errors into one |
| Custom type with `Unwrap() error` | Yes | `type QueryError struct { Op string; Err error }` | Rich, structured errors callers can inspect with `errors.As` |
| `fmt.Errorf(..., %v)` | No | `fmt.Errorf("failed: %v", err)` | Simple message formatting only; original error is not unwrappable |

### `fmt.Errorf` with `%w` (most common)

```go
err := fmt.Errorf("read %s failed: %w", path, err)
```

### `errors.Join` (Go 1.20+)

```go
err := errors.Join(err1, err2)
```

Returns `nil` if all inputs are `nil`. `errors.Is` and `errors.As` check against any error in the joined set.

### Custom error type with `Unwrap()`

```go
type QueryError struct {
    Op  string
    Err error
}

func (e *QueryError) Error() string { return fmt.Sprintf("%s: %v", e.Op, e.Err) }
func (e *QueryError) Unwrap() error { return e.Err }
```

## Run it

```bash
go run main.go
```

A dummy `database/sql` driver is registered so `Query` succeeds but `Close` fails. Output:

```
getBalance1: balance=0, err=<nil>         // close error silently discarded
getBalance2: balance=0, err=<nil>         // close error explicitly ignored
getBalance3: balance=0, err=close failed  // close error captured and returned
```

This shows the difference between the three approaches: only `getBalance3` surfaces the close error to the caller.

## Key takeaway

Never discard errors from deferred cleanup calls with `defer f()` or `defer func() { _ = f() }()`. Capture the error and decide whether to log it, overwrite the return error, or wrap both together.
