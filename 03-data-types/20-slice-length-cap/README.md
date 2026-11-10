# Slice Length vs Capacity

This is sample code from the book *100 Go Mistakes* (mistake #20 in the "Data Types" chapter).

It demonstrates how Go slices really work under the hood — specifically the difference between **length** and **capacity**, and how **slices share backing arrays**, which leads to subtle aliasing bugs.

## What the code shows

### 1. `len` vs `cap`

`s := make([]int, 3, 6)` creates a slice with `len=3` (elements you can read/write) and `cap=6` (room in the underlying array before a reallocation is forced).

### 2. `append` may or may not copy

- `append(s, 2)` fits within the spare capacity → `len` grows to 4, `cap` stays 6.
- `append(s, 3, 4, 5)` exceeds capacity → Go allocates a new, larger backing array (capacity roughly doubles, 6 → 12) and copies the old data over.

### 3. Slices are descriptors over a shared array

`s2 := s1[1:3]` creates a *new slice header* (pointer, length, capacity) that points into `s1`'s same backing array. Its capacity is `cap(s1) - 1 = 5`.

Because they share memory, writing through one slice shows up in the other:

```go
s1[1] = 1
print(s2) // len=2, cap=5: [1 0]  — the change is visible through s2
```

### 4. The aliasing gotcha

Appending through `s2` reuses the shared backing array as long as there is spare capacity:

```go
s2 = append(s2, 2)   // writes into s1's array at index 3
print(s1)            // len=3, cap=6: [0 1 0] — s1's len is 3, so s1[3] is hidden
```

But once `s2` outgrows its capacity, Go allocates a fresh array for `s2`. After that point, `s1` and `s2` no longer share memory, so further appends to `s2` leave `s1` untouched:

```go
s2 = append(s2, 5)   // forces reallocation; s2 gets a new array (cap 10)
print(s1)            // len=3, cap=6: [0 1 0] — unchanged
print(s2)            // len=6, cap=10: [1 0 2 3 4 5] — on its own array now
```

## The mistake it warns against

Assuming that a slice owns its data, or assuming that `append` always gives you an independent copy.

Two slices can alias the same array, so a modification (or an `append` that fits in spare capacity) through one slice can silently corrupt another — until a reallocation breaks the sharing.

## The fix

When you need an independent copy, don't rely on slicing. Explicitly copy instead:

```go
s2 := append([]int{}, s1...)
// or
s2 := make([]int, len(s1))
copy(s2, s1)
```

## Output

```
len=3, cap=6: [0 0 0]
len=3, cap=6: [0 1 0]
len=4, cap=6: [0 1 0 2]
len=7, cap=12: [0 1 0 2 3 4 5]
len=2, cap=5: [1 0]
len=3, cap=6: [0 1 0]
len=3, cap=5: [1 0 2]
len=3, cap=6: [0 1 0]
len=6, cap=10: [1 0 2 3 4 5]
```
