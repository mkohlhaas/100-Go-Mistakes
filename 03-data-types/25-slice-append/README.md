# Slice append aliasing — Mistake #25

This example (from *100 Go Mistakes*, mistake #25) demonstrates a classic Go slice gotcha:

> **`append` can silently mutate a *different* slice that shares the same underlying array — when the slice you append to has spare capacity.**

## How to run

```bash
cd /home/schmidh/Gitrepos/100-go-mistakes/03-data-types/25-slice-append
go run main.go
```

Output:

```
[1 2 10]
[1 2 3]
[1 2 3]
```

## What's happening

A slice is a small header (pointer, length, capacity) over a shared backing array. Slicing (`s[:2]`) does **not** copy data — it just creates a new header pointing at the same array.

### `listing1()` — the bug

```go
s := []int{1, 2, 3}   // len 3, cap 3, backing array [1 2 3]
f(s[:2])              // same array, len 2, cap 3 → room to grow
fmt.Println(s)        // prints [1 2 10] ← the 3 is gone!
```

Inside `f`, `append(s[:2], 10)` finds spare capacity (`cap 3 > len 2`), so it writes `10` at index 2 of the **shared** backing array without allocating. `s` still has length 3 over that same array, so it now reads `[1 2 10]`. The original `3` was silently overwritten.

### `listing2()` — fix via explicit copy

```go
s := []int{1, 2, 3}
sCopy := make([]int, 2)   // len 2, cap 2 — no spare room
copy(sCopy, s)             // [1 2]
f(sCopy)                   // append inside f discards its result; sCopy stays [1 2]
result := append(sCopy, s[2])
fmt.Println(result)        // prints [1 2 3] ✓
```

Step by step:

1. `sCopy := make([]int, 2)` creates a slice with `len 2, cap 2` backed by its own fresh array `[0 0]`.
2. `copy(sCopy, s)` copies `1, 2` into it, so `sCopy` is now `[1 2]`. `sCopy` and `s` no longer share an array.
3. `f(sCopy)` calls `append(sCopy, 10)`. Because `cap == len` (no spare room), Go allocates a **new** backing array `[1 2 10]`. That result is assigned to `_` and thrown away. `sCopy`'s header still points to its original `[1 2]` array — it was **not** modified.
4. `result := append(sCopy, s[2])` appends `s[2]`, which is `3` (from the original `s`). Again `cap == len`, so a fresh array is allocated: `[1 2 3]`.

`make([]int, 2)` gives `cap == len`, so any `append` on `sCopy` must allocate a new backing array. `s`'s data is never touched.

### `listing3()` — fix via three-index slice

```go
f(s[:2:2])   // cap capped at 2 (no spare room)
// prints [1 2 3] ✓
```

The three-index form `s[low:high:max]` limits capacity to `max`, so `append` can't reuse `s`'s spare capacity and must allocate fresh.

## The takeaway

Never assume `append` is safe to call on a slice you sliced from another slice. If there's spare capacity, `append` writes into the shared array and can corrupt the original slice's contents. To be safe, either:

1. Copy into a fresh slice made with `make([]T, n)` (`cap == len`), or
2. Use the three-index expression `s[:i:i]` to cap capacity before appending.
