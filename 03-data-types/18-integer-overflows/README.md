# Integer Overflow — Mistake #18

This directory contains sample code from *100 Go Mistakes* demonstrating how to **detect and prevent integer overflow in Go**.

## The problem

Go does **not** raise an error when integer math overflows. The result silently wraps around instead:

- For unsigned types, wrap-around is defined modulo 2ⁿ.
- For signed types, wrap-around effectively happens in practice.

So code like `counter + 1` is unsafe when `counter` is already at the type's maximum — it wraps to the minimum rather than failing loudly.

## What each function shows

| Function | Technique demonstrated |
|----------|------------------------|
| `Inc32` | Guard an increment: panic if `counter == math.MaxInt32` instead of letting it wrap to `math.MinInt32`. |
| `IncInt` | Same guard for the platform-dependent `int` type, using `math.MaxInt`. |
| `IncUint` | Same guard for `uint`, using `math.MaxUint`. |
| `AddInt` | Check `a > math.MaxInt-b` before adding. The comparison is rearranged so the check itself cannot overflow. |
| `MultiplyInt` | Multiplication guard with edge cases: <br>• Zero and multiply-by-1 short-circuits <br>• `math.MinInt` special case (negating `MinInt` does not fit, e.g. `MinInt * -1` overflows) <br>• Division check `result/b != a` — if division does not reproduce `a`, the multiplication overflowed. |

## Key takeaways

1. **Check before you compute.** Validate operands against the type's limits (`math.MaxInt`, `math.MaxUint`, `math.MaxInt32`, etc.) before adding or multiplying.
2. **Rearrange comparisons.** `a + b > MaxInt` can itself overflow; write it as `a > MaxInt - b` instead.
3. **Watch `math.MinInt`.** It is the asymmetric edge case of signed integers: `-MinInt` is out of range, so `MinInt * -1` overflows.
4. **Verify multiplication via division.** If `(a * b) / b != a`, the product overflowed (provided `b ≠ 0`).

## Bottom line

The language will not protect you from silent integer overflow. Use defensive patterns like those in `main.go` to validate operands before performing arithmetic.

## Go vs. Rust: no built-in checked arithmetic

Unlike Rust, Go does **not** provide checked arithmetic methods on its integer types. Rust lets you opt in explicitly:

```rust
a.checked_add(b)   // Option<result>, None on overflow
a.overflowing_add(b) // (result, overflowed_bool)
a.saturating_add(b)  // clamps at max/min
a.wrapping_add(b)    // wraps silently
```

Go has nothing like that. The `+`, `-`, and `*` operators just compute silently:

- **Unsigned**: wrap-around is defined modulo 2ⁿ.
- **Signed**: overflow does not panic or return an error; it effectively wraps (the spec allows it, and in practice it wraps in two's complement).

So you must guard manually, exactly as the sample code does:

```go
if a > math.MaxInt-b {
    panic("int overflow")
}
return a + b
```

### The one exception: `math/bits`

Go's `math/bits` package does have low-level helpers that return a carry/overflow bit:

```go
import "math/bits"

sum, carry := bits.Add(uint(a), uint(b), 0)
if carry != 0 {
    // overflow
}

hi, lo := bits.Mul(uint(a), uint(b))
if hi != 0 {
    // multiplication overflowed beyond uint range
}
```

These are verbose and operate on `uint`, so most Go code uses manual range checks like those in `main.go` instead.

**Bottom line:** Rust makes overflow handling explicit and opt-in via methods; Go leaves it to you to check before you compute.