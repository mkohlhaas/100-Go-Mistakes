# Range loop argument evaluation — arrays

Example from *100 Go Mistakes and How to Avoid Them*, mistake #31:
the range loop argument is evaluated **once** before the loop starts,
and for arrays it is evaluated as a **copy**. Modifying the array while
iterating does not change what the iteration variable yields — unless
you range over a pointer to the array.

## Files

- `main.go` — three listings showing the behavior:
  - `listing1()` — `for i, v := range a` then mutates `a[2]`. Prints `2`
    (the original value, because `a` was copied before the loop).
  - `listing2()` — `for i := range a` then reads `a[2]` directly inside
    the body. Prints `10` (the mutation is visible because the element
    is read fresh from the array each time).
  - `listing3()` — `for i, v := range &a` then mutates `a[2]`. Prints `10`
    (ranging over `&a` walks the array through the pointer, so mutations
    are reflected in `v`).

## Run it

```sh
go run main.go
```

Output:

```
2
[0 1 10]
10
[0 1 10]
10
[0 1 10]
```

Each listing also prints the final state of `a` (`[0 1 10]`), confirming the array itself was mutated in all three cases — only what the loop variable sees differs.

## Key takeaway

The range expression is evaluated exactly once, before iteration begins.
For arrays (and slices), that evaluation copies the data, so the loop
variable `v` holds values from the original copy even if you mutate the
underlying array during iteration. If you need to see mutations inside
the loop, read directly from the array/slice by index (like `listing2`)
or range over a pointer to the array (like `listing3`).
