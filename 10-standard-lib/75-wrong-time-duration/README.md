# Mistake 75 — Wrong time duration

Example from *100 Go Mistakes and How to Avoid Them*, Chapter 10 (Standard Library).

## The mistake

`time.Duration` is an `int64` representing **nanoseconds**. Passing a bare
integer to a function that expects a `time.Duration` does not mean milliseconds
or seconds — it means nanoseconds.

```go
// BUG: 1000 nanoseconds = 1 microsecond, NOT 1 second
ticker := time.NewTicker(1000)
```

This ticker fires every microsecond instead of every second — 1,000× faster
than intended.

## The fix

Always express durations explicitly using the `time` package's duration
constants, or multiply a plain number by one of those constants:

```go
// GOOD: explicit duration constants
ticker := time.NewTicker(time.Second)

// GOOD: plain numbers must be multiplied by a duration constant
ticker := time.NewTicker(1000 * time.Millisecond)
```

Available constants: `time.Nanosecond`, `time.Microsecond`, `time.Millisecond`,
`time.Second`, `time.Minute`, `time.Hour`.

## Run it

`listing1` shows the buggy bare-integer version; `listing2` shows the fixed
version with an explicit duration. Switch between them in `main()`:

```bash
go run main.go
```

## Key takeaway

Never pass a bare integer where a `time.Duration` is expected. Use `time`
duration constants, or multiply by one, so the unit is explicit and the intent
is clear.
