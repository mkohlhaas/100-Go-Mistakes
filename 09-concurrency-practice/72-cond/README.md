# 72-cond — `sync.Cond` condition variables

Practice for Chapter 9 (Concurrency), item 72 of *100 Go Mistakes and How to Avoid Them*:
**wait on a condition with `sync.Cond` instead of polling**.

## What it does

Several "listener" goroutines wait until a shared `balance` reaches a goal amount, while one "updater" goroutine increments the balance every second. Three implementations are compared:

### Listing 1 — Busy polling with `sync.Mutex` ❌

Listeners poll: lock, check `balance < goal`, unlock, sleep briefly, re-lock. It eventually reaches the goal, but burns CPU and is slow to react. (The original `sync.RWMutex` version is worse: readers starve the writer, so the balance never increments and the program hangs.) This is the anti-pattern `sync.Cond` exists to replace.

### Listing 2 — Channel notification

The updater sends each new balance on a channel; listeners `range` over it and exit when the goal is met. Simple and correct, but every listener sees every update, which can be wasteful when many listeners wait for different goals.

### Listing 3 — `sync.Cond` ✅

The idiomatic way to block until a condition becomes true:

- Listeners lock `cond.L`, call `cond.Wait()` inside a `for` loop while `balance < goal`, then unlock and report. `Wait()` atomically unlocks the mutex and suspends the goroutine — no CPU is wasted while waiting.
- The updater locks the mutex, increments `balance`, unlocks, and calls `cond.Broadcast()` to wake all waiters.

## How `sync.Cond` works

A `*sync.Cond` is a rendezvous point around a lock (`sync.Locker`):

1. A waiter calls `cond.Wait()` while holding `cond.L`. It atomically releases the lock and sleeps; on wakeup it re-acquires the lock before returning.
2. A notifier changes the shared state, then calls `cond.Signal()` (wake one waiter) or `cond.Broadcast()` (wake all waiters).

### Rules to follow

- **Always hold `cond.L` when calling `Wait()`** — it unlocks and suspends atomically.
- **Always wait in a `for` loop re-checking the condition**, never a bare `if`. Wakeups can be spurious, and another goroutine may have already satisfied (or invalidated) the condition while you were asleep.
- **Call `Signal()`/`Broadcast()` after the state change is visible** (typically while still holding the lock, or right after releasing it) to avoid missed wakeups.

## Running

```bash
cd /home/schmidh/Gitrepos/100-Go-Mistakes/09-concurrency-practice/72-cond
go run main.go
```

`main.go` runs `listing3` (the recommended `sync.Cond` version) by default, with goals of `$2` and `$5` — expect output within a few seconds. To try a different listing, change the call in `main()` to `listing1()` or `listing2()`.

## Takeaway

Don't poll (Listing 1). Use channels for simple fan-out (Listing 2). Use `sync.Cond` when goroutines need to sleep efficiently until a shared condition becomes true (Listing 3).