# Nil Receiver / Typed Nil Error

This example demonstrates a common Go mistake: returning a **typed nil pointer** where an `error` interface is expected.

## The mistake

`Validate1()` builds a `*MultiError` and returns it directly:

```go
var m *MultiError
// ... possibly assign m = &MultiError{} ...
return m
```

When no error occurred, `m` is a **typed nil**: its type is `*MultiError`, but its value is `nil`.

In Go, an `error` is an interface. You can think of an interface as a box with two slots inside:

| Slot | What it holds |
|------|---------------|
| Type | The concrete type (e.g. `*MultiError`) |
| Value | The actual data (e.g. a pointer, or `nil`) |

An interface is `nil` only when **both** slots are empty. When you return a typed nil pointer as an `error`, Go fills the type slot with `*MultiError` and the value slot with `nil`. Because the type slot is filled, the interface itself is **not nil** — it's a non-nil interface wrapping a nil value.

Consequently, every caller sees a non-nil error:

```go
if err := customer.Validate1(); err != nil {   // SA4023: always true
    log.Fatalf("customer is invalid: %v", err)
}
```

Even for a perfectly valid customer, `err` is a non-nil interface wrapping a nil `*MultiError`. So `err != nil` is always true — that's the **SA4023: this comparison is always true** warning from staticcheck, because `Validate1()` can never produce a genuine `nil` error.

## The fix

`Validate2()` shows the correct pattern: explicitly return a bare `nil` when there is nothing to report.

```go
if m != nil {
    return m
}
return nil
```

A bare `nil` puts nothing in either slot of the interface box, so `err` is genuinely `nil` and `err != nil` correctly becomes `false` for a valid customer.

In this example, `main()` calls `Validate2()`, so the program runs without fatal logging. `Validate1()` is left unchanged to show what *not* to do.

## Key takeaway

Never return a typed nil pointer (or any typed nil value) as an `error` interface. When in doubt, return the untyped literal `nil`:

```go
return nil   // good: untyped nil → interface is truly nil
```

Returning `&SomeType{}` wrapped in an interface, even when the pointer is nil, makes the interface non-nil and breaks every `err != nil` check downstream.
