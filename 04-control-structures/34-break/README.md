# `break`

Example from **"100 Go Mistakes"** — Chapter 4: Control Structures, item 34.

## The mistake

In Go, a bare `break` inside a `switch` or `select` only breaks out of that `switch`/`select` — **not** out of any enclosing `for` loop. To exit an enclosing loop, you need a labeled break (`break <label>`).

## What it shows

- **listing1** — `break` inside the `switch` only exits the `switch`. The `for` loop keeps running, so all of `0 1 2 3 4` is printed.
- **listing2** — `break loop` breaks out of the labeled `for` loop. When `i == 2`, both the `switch` and the loop exit, printing `0 1 2`.
- **listing3** — `break` inside the `select` only exits the `select`. When the context is canceled, the loop spins on the now-ready `ctx.Done()` case and never returns.
- **listing4** — `break loop` breaks out of the labeled `for` loop. When the context is canceled, the loop and the function return.

## Run it

```sh
go run main.go
```

## Fix

Use a labeled break when you want to exit the enclosing loop rather than just the `switch`/`select`:

```go
loop:
	for {
		select {
		case <-ctx.Done():
			break loop // exits the for loop
		}
	}
```

## Learn more

- [Go spec — Break statements](https://go.dev/ref/spec#Break_statements)
- [Effective Go — Break and continue](https://go.dev/doc/effective_go#break-continue)
