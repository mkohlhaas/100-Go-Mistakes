# 17 — Octal literals

## What the code demonstrates

Go supports octal (base-8) integer literals using the `0o` (or `0O`) prefix:

```go
sum := 100 + 0o10
fmt.Println(sum) // 108
```

`0o10` is an octal number, not decimal 10:

```
0o10 = 1×8¹ + 0×8⁰ = 8 (decimal)
```

So `100 + 0o10` equals `100 + 8`, which is **108**, not `110`.

## The mistake

The pitfall is assuming `0o10` means decimal ten. A reader unfamiliar with octal notation might expect the sum to be `110`, but it is actually `108`.

Go 1.13 introduced the explicit `0o`/`0O` prefix (replacing the older C-style leading zero, e.g. `010`) to make octal intent clearer, but the prefix itself can still trip people up who don't realize it changes the base.

**Takeaway:** literals prefixed with `0o` are interpreted in base 8, not base 10.
