# 23 — Checking slice empty

This sample demonstrates a common Go mistake: **checking whether a slice is empty by comparing it to `nil` instead of checking its length.**

## The code

```go
func handleOperations1(id string) {
	operations := getOperations(id)
	if operations != nil {   // ← mistake: checks nil-ness, not emptiness
		handle(operations)
	}
}

func handleOperations2(id string) {
	operations := getOperations(id)
	if len(operations) != 0 {   // ← correct: checks if there are any elements
		handle(operations)
	}
}

func getOperations(id string) []float32 {
	operations := make([]float32, 0)

	if id == "" {
		return operations
	}

	// Add elements to operations

	return operations
}

func handle(operations []float32) {}
```

## Why it matters

In Go, a **nil slice** and an **empty non-nil slice** (`make([]float32, 0)`) are different values, but both have length 0. Crucially, `len(nil)` is `0` in Go.

`getOperations` returns `make([]float32, 0)` — a non-nil but empty slice — when `id == ""`:

- **`handleOperations1`** (`operations != nil`): evaluates to `true` for that empty non-nil slice, so it calls `handle` with an empty slice even though there is nothing to process.
- **`handleOperations2`** (`len(operations) != 0`): evaluates to `false` for both nil and empty slices, so `handle` is only called when there are actual elements.

## The rule

To check if a slice has elements, use `len(slice) == 0` (empty) or `len(slice) != 0` (non-empty).

Do **not** use `slice == nil` or `slice != nil` for emptiness checks — `len` already treats nil and empty slices identically, and an empty non-nil slice will fool a nil comparison.
