# `08-any` — Avoid `any` when you know the type

This sample comes from chapter **08** ("any") of *100 Go Mistakes*.

## What it demonstrates

Go's `any` type (an alias for `interface{}`) can hold a value of any type:

```go
var i any
i = 42
i = "foo"
i = struct{ s string }{s: "bar"}
i = f
```

That looks flexible, but it is **type-unsafe**: nothing stops you from assigning the wrong kind of value, and nothing tells you what `i` actually holds until you do a runtime type assertion.

The real mistake appears in `store/before.go`, where a `Store` uses `any` everywhere:

```go
func (s *Store) Get(id string) (any, error)
func (s *Store) Set(id string, v any) error
```

Because everything is `any`, you can `Set` a `Customer` and `Get` it back expecting a `Contract`. The compiler won't catch the mismatch — it fails at runtime or silently produces wrong behavior.

## The fix

`store/after.go` replaces the `any` API with concrete, type-specific methods:

```go
func (s *Store) GetContract(id string) (Contract, error)
func (s *Store) SetCustomer(id string, customer Customer) error
```

Now type mismatches are caught at compile time, no type assertions are needed, and the code documents itself.

## Takeaway

Don't reach for `any` just because it's convenient. If you know the type, use it. Reserve `any` for genuinely heterogeneous data (e.g., logging, unstructured JSON). With Go 1.18+, prefer generics when you need truly type-safe, reusable code over `any`.

## Files

- `main.go` — shows `any` accepting unrelated value types
- `store/before.go` — the mistake: `any`-based store API
- `store/after.go` — the fix: typed store API