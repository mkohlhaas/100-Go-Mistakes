# Mistake #24 — Assuming `copy` Allocates the Destination Slice

This code demonstrates a common Go mistake: calling `copy(dst, src)` when `dst` is a nil (or empty) slice, expecting it to allocate space and copy all elements. It doesn't.

## The problem (`bad()`)

```go
src := []int{0, 1, 2}
var dst []int          // dst is nil, len 0
copy(dst, src)         // copies nothing!
fmt.Println(dst)       // prints []
```

`var dst []int` creates a nil slice with length 0. Go's builtin `copy` copies only `min(len(dst), len(src))` elements — here, `min(0, 3) = 0`. So nothing is copied and `dst` stays empty; the elements are silently dropped.

## The fix (`correct()`)

```go
src := []int{0, 1, 2}
dst := make([]int, len(src))  // pre-allocate dst with same length as src
copy(dst, src)                // copies all 3 elements
fmt.Println(dst)              // prints [0 1 2]
```

By allocating `dst` first with `make`, it has room for all elements, so `copy` works as expected.

## Key takeaway

Go's `copy(dst, src)` **never allocates memory** for `dst` and **never expands it**. It only copies into existing capacity, up to the destination's length. If `dst` is shorter than `src`, the tail of `src` is silently truncated.

To properly duplicate a slice, you must allocate the destination first — e.g. `make([]int, len(src))` — or use `append(dst[:0:0], src...)`.