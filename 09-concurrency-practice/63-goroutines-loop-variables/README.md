# Goroutines and Loop Variables

This project demonstrates the classic Go pitfall of capturing a loop variable inside a goroutine, and shows two common ways to fix it.

## The problem

`listing1()` spawns a goroutine per element of a slice:

```go
for _, i := range s {
    go func() {
        fmt.Print(i)
    }()
}
```

Before Go 1.22, the loop variable `i` was a single variable reused across all iterations. Every goroutine captured a reference to that same variable — not the value it held at creation time. Because goroutines are scheduled lazily, by the time they actually run the loop has usually finished, and `i` holds its final value (`3`). The program typically prints `333` instead of `123`.

**Go 1.22+ note:** The language spec changed so each loop iteration gets its own variable instance. On Go 1.22 or later `listing1()` will print `123` (order not guaranteed). The examples here target pre-1.22 semantics, which is the form still found in much existing code and the original source of this mistake.

## The fixes

### Fix 1 — copy into a local variable (`listing2`)

Introduce a new variable scoped to each iteration:

```go
for _, i := range s {
    val := i
    go func() {
        fmt.Print(val)
    }()
}
```

Each goroutine now captures its own `val`, so it prints the value correct to its iteration.

### Fix 2 — pass the value as an argument (`listing3`)

Pass `i` as a parameter to the goroutine function:

```go
for _, i := range s {
    go func(val int) {
        fmt.Print(val)
    }(i)
}
```

Arguments are evaluated when the goroutine is created, so each call receives the value of `i` at that moment.

## Running the code

```bash
go run main.go
```

`main()` calls `listing1()`, `listing2()`, and `listing3()` in turn, each prefixed with a label:

```
listing1: 321
listing2: 123
listing3: 312
```

Each listing uses a `sync.WaitGroup` so `main()` waits for all goroutines to finish before printing, guaranteeing output appears. Goroutine execution order isn't guaranteed, so the digits may appear in a different order on each run.

On Go 1.22+ (this environment runs 1.27), `listing1()` no longer exhibits the bug — each loop iteration gets its own `i` — so all three listings print each value once. To observe the pre-1.22 buggy behavior (all goroutines printing the final value, e.g. `333`), run with an older toolchain.

## Key takeaway

Never capture a loop variable directly in a closure/goroutine (pre-Go 1.22). Either copy it into a per-iteration local, or pass it as a function argument.

## Reference

This example is based on mistake #63 from *100 Go Mistakes and How to Avoid Them* by Teiva Harsanyi.