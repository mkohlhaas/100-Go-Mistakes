# Defer argument evaluation

This example accompanies chapter **47 — Defer evaluation** of *100 Go Mistakes and How to Avoid Them*.

## The mistake

Go evaluates a `defer` statement's arguments **when the `defer` executes**, not when the deferred call runs. Passing a value-type argument means the deferred call snapshots the variable's current (zero) value:

```go
// BUG: status is captured as "" — the value at defer time
var status string
defer notify(status)      // prints "notify: "
defer incrementCounter(status) // prints "increment: "

if err := foo(); err != nil {
    status = StatusErrorFoo
    return err
}
...
status = StatusSuccess
return nil
```

`notify` and `incrementCounter` always print the empty string, even though `status` is updated before the function returns. The deferred calls captured a copy of `status` at the moment the `defer` statements ran.

## The fix

Either pass a pointer so the deferred call reads the variable's **current** value at execution time:

```go
// GOOD: pointer read at execution time reflects the latest value
var status string
defer notifyPtr(&status)
defer incrementCounterPtr(&status)
```

...or wrap the deferred calls in a closure that captures the variable by reference:

```go
// GOOD: closure captures status, read at execution time
var status string
defer func() {
    notify(status)
    incrementCounter(status)
}()
```

Both approaches observe `status` after it has been assigned (`StatusErrorFoo`, `StatusErrorBar`, or `StatusSuccess`).

## What this example shows

* `f1` — **bug**: value arguments snapshot `""` at defer time.
* `f2` — **fix**: pointers let the deferred calls read the final value.
* `f3` — **fix**: a closure captures `status` by reference.

## Running

```bash
go run main.go
go vet ./...
go test ./...
```

Expected output from `f2` and `f3` is `notify: success` / `increment: success`; `f1` prints empty values.

```bash
# Output:
increment:
notify:
increment: success
notify: success
notify: success
increment: success
```
