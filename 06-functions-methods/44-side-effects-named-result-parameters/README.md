# Side effects with named result parameters

This example accompanies chapter **44 — Side effects with named result parameters** of
*100 Go Mistakes and How to Avoid Them*.

## The mistake

Go lets you name a function's return values in its signature:

```go
func (l loc) getCoordinates1(ctx context.Context, address string) (
	lat, lng float32, err error) {
```

The names `lat`, `lng`, and `err` are declared in the function's scope and
initialized to their zero values. That makes them ordinary variables: anything
assigned to them persists until the function returns, and a bare `return`
returns their current values.

The mistake is treating a named result as if it were an output-only slot that
only an explicit `return x, y, err` can touch. In `getCoordinates1`:

```go
	if ctx.Err() != nil {
		return 0, 0, err
	}
```

`ctx.Err()` is called, but its result is never stored anywhere. The named
result `err` is still its zero value (`nil`), so `return 0, 0, err` returns a
nil error even though the context is cancelled. The context error is silently
swallowed — a side effect on the named result that never happened.

Run the example:

```bash
go run main.go
```

Output:

```
getCoordinates1: lat=0 lng=0 err=<nil> (want a context error, got true)
getCoordinates2: lat=0 lng=0 err=context canceled
getCoordinates3: lat=0 lng=0 err=context canceled
```

Only `getCoordinates1` loses the error.

## The three variants

| Method | What it does | Result |
|--------|--------------|--------|
| `getCoordinates1` | calls `ctx.Err()` but discards the result; returns the named `err` (still `nil`) | **BUG** — error swallowed |
| `getCoordinates2` | `if err := ctx.Err(); err != nil` declares a new `err` that shadows the named result, then returns it explicitly | Works, but relies on shadowing; fragile and easy to misread |
| `getCoordinates3` | `if err = ctx.Err(); err != nil` assigns to the named result (no `:=`, so no new variable), then a bare `return` returns it | Correct and clear |

`getCoordinates2` happens to work because the `:=` inside the `if` creates a
new variable named `err` that shadows the named result, and the explicit
`return 0, 0, err` returns that shadowed variable. But if you later change the
`return` to a bare `return`, or move the check, the named result (still `nil`)
is returned instead. Don't rely on it.

## The fix

Assign to the named result with `=` (not `:=`), so you're updating the
existing named result rather than declaring a shadowing local:

```go
	if err = ctx.Err(); err != nil {
		return
	}

	// Get and return coordinates
	return 0, 0, nil
```

Or simply return the error explicitly:

```go
	if err := ctx.Err(); err != nil {
		return 0, 0, err
	}
```

## Rules of thumb

1. Named results are variables, not just documentation — reading or writing
   them has real effects.
2. A bare `return` returns the *current* values of the named results. If you
   never assigned to one, you return its zero value.
3. Avoid `:=` with the same name as a named result inside the function; it
   creates a shadowing local and the named result keeps its old (zero) value.
4. When in doubt, return values explicitly rather than with a bare `return`.

## Running

```bash
go run main.go
go vet ./...