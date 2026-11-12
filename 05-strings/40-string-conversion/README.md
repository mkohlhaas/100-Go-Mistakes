# Mistake #40 — Round-Tripping Between `string` and `[]byte`

This code demonstrates why converting back and forth between `string` and `[]byte` is wasteful, and how to avoid it by staying in `[]byte` land.

## The code

Two ways to read from an `io.Reader` and trim surrounding whitespace:

* `getBytes1` — `[]byte` → `string` → `strings.TrimSpace` → `string` → `[]byte` (the mistake: two conversions)
* `getBytes2` — `[]byte` → `bytes.TrimSpace` → `[]byte` (no conversions)

## Why `getBytes1` is wasteful

In Go, converting `[]byte` to `string` copies the bytes (strings are immutable and need their own backing array), and converting `string` back to `[]byte` copies again. `getBytes1` pays that copy cost twice on every call:

```
[]byte --copy--> string --copy--> []byte
```

`getBytes2` trims directly on the `[]byte` with `bytes.TrimSpace`, which avoids the round-trip and can reuse the existing allocation (it just re-slices the input).

The same family of functions exists for `[]byte`: `bytes.Trim`, `bytes.TrimLeft`, `bytes.TrimRight`, `bytes.TrimPrefix`, `bytes.TrimSuffix`, and `bytes.TrimSpace` mirror their `strings` counterparts but operate on slices without any string conversion.

## Run it

```bash
go run main.go
```

Example output:

```
input: "   hello, world!   "
getBytes1: "hello, world!"
getBytes2: "hello, world!"
```

Both produce the same result, but `getBytes2` avoids two full allocations per call.

## Key takeaway

When processing `[]byte`, keep it as `[]byte` as long as possible. Convert to `string` only when you actually need one (for example, to use it as a map key or to call an API that requires strings), and avoid round-trips like `[]byte → string → []byte`.

## References

- [`bytes.TrimSpace`](https://pkg.go.dev/bytes#TrimSpace)
- [`strings.TrimSpace`](https://pkg.go.dev/strings#TrimSpace)
- [100 Go Mistakes — mistake #40: string conversion](https://100go.co)
