# Floating-Point Mistakes

This example demonstrates three common pitfalls when working with floating-point numbers in Go.

## 1. Limited precision (`float32`)

```go
var n float32 = 1.0001
fmt.Println(n * n) // 1.0002
```

`1.0001` cannot be represented exactly in binary floating point, so squaring it produces rounding error. The result is an approximation, not the mathematically exact `1.00020001`. `float32` is less precise than `float64`, which makes this worse.

## 2. Division by zero yields `Inf`/`NaN`, not a panic

```go
var a float64
positiveInf := 1 / a      // +Inf
negativeInf := -1 / a     // -Inf
nan := a / a              // NaN
```

Unlike integer division by zero (which panics at runtime), floating-point division by zero follows IEEE 754 and silently yields `+Inf`, `-Inf`, or `NaN`. Code that assumes a panic — or that compares `Inf` or `NaN` with `==` — will behave unexpectedly (notably, `NaN != NaN`).

## 3. Order of operations changes the result

```go
fmt.Println(f1(100), f2(100)) // 10100.00999999993 10100.01
```

Both functions mathematically equal `10000 + 100×1.0001 = 10100.01`, but they print different values:

- `f1` adds `1.0001` repeatedly to an accumulator already at `10000`. Each small addition gets rounded relative to the large magnitude, so error accumulates → `10100.00999999993`.
- `f2` accumulates the small values first, then adds `10000` at the end → `10100.01`.

**Takeaway:** floating-point addition is not associative. When summing values of very different magnitudes, add the small values together first, then combine with the large value — and avoid comparing floats with `==`.
