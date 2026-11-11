# Comparing Values (Mistake #29)

This example demonstrates the pitfalls of comparing values in Go, particularly with
structs that contain non-comparable fields such as slices.

## What it shows

- **Comparable structs** (`customer1`, only string fields) can be compared with `==`.
- **Non-comparable structs** (`customer2` containing a `[]float64`) cannot be compared
  with `==` — the code won't compile (`invalid operation: struct containing []float64 cannot be compared`).
- **`any` / `interface{}` comparisons**: two `any` values holding the same comparable
  concrete value compare equal; two `any` values holding a non-comparable concrete
  type (e.g. a struct with a slice) **panic at runtime** instead of returning `false`
  (shown commented out to avoid the panic).
- **`reflect.DeepEqual`** can deep-compare non-comparable values safely.
- A custom `equal` method provides an explicit, field-by-field equality check.

## Run it

```bash
go run main.go
```

Expected output:

```
true
true
true
true
false
```

## Key takeaway

Don't rely on `==` for structs that may contain slices, maps, or functions. Use
`reflect.DeepEqual` or write an explicit equality method (like `customer2.equal`)
instead.
