# 49 – Error Wrapping

Example from _100 Go Mistakes_ demonstrating how to wrap errors correctly in Go so they can still be inspected with `errors.Is` and `errors.Unwrap`.

## The problem

When wrapping an error you must use the `%w` verb in `fmt.Errorf`. Using `%v` (or a custom wrapper type that doesn't implement `Unwrap`) loses the original error, so callers can no longer unwrap or match it with `errors.Is`.

## Code walkthrough

`main.go` contains four listings:

- **listing1** – returns the error unwrapped. Simple, but adds no context.
- **listing2** – wraps the error in a custom `BarError` struct using string concatenation. The original error is preserved as a field, but because there is no `Unwrap() error` method, `errors.Is` and `errors.Unwrap` cannot reach it.
- **listing3** – **correct.** Uses `fmt.Errorf("bar failed: %w", err)`. The `%w` verb wraps the error and makes it accessible via `errors.Unwrap`, so `errors.Is(err, barError{})` works.
- **listing4** – **incorrect.** Uses `fmt.Errorf("bar failed: %v", err)`. The `%v` verb formats the error into a plain string; the original error is lost and cannot be unwrapped.
- **listing5** – **longer chain.** Wraps `listing3`'s error with another `%w`, producing a chain `listing5 -> listing3 -> bar`. `errors.Unwrap` can walk each layer, and `errors.Is` still matches the root `barError{}` at the bottom of the chain.

## Walking the error chain

You can traverse a wrapped error with a `for` loop using `errors.Unwrap`. Each call returns the next error in the chain, or `nil` at the end:

```go
for err := err5; err != nil; err = errors.Unwrap(err) {
    fmt.Printf("  %v\n", err)
}
```

`errors.Is` and `errors.As` walk this same chain internally, which is why they can still reach the root error no matter how many layers deep it is.

## Run it

```bash
go run main.go
```

## Quintessence

Preserve the error chain. Wrap errors with context, but never lose the original error.

Do as in `listing3()`:

```go
err := bar()
if err != nil {
    return fmt.Errorf("bar failed: %w", err)
}
```

## Key takeaway

Always use `%w` (not `%v`) when wrapping errors with `fmt.Errorf`, and ensure custom wrapper types implement `Unwrap() error`, so the error chain remains inspectable.
