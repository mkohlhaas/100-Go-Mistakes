# 85. Table-Driven Tests

Example from *100 Go Mistakes and How to Avoid Them* — Chapter 11: Testing.

## What this example does

`main.go` provides a small utility function, `removeNewLineSuffixes`, that recursively strips trailing newline characters (`\n` and `\r\n`) from a string:

```go
func removeNewLineSuffixes(s string) string
```

Behavior:
- `""` → `""`
- `"a\r\n"` → `"a"`
- `"a\n"` → `"a"`
- `"a\n\n\n"` → `"a"`
- `"a"` → `"a"` (no trailing newline, returned unchanged)

## The point: table-driven tests

`main_test.go` shows two ways to test the same behavior:

1. **Individual test functions** (`TestRemoveNewLineSuffix_Empty`, `TestRemoveNewLineSuffix_EndingWithNewLine`, etc.) — one function per case. This is repetitive and hard to maintain; adding a case means adding a whole new test function.
2. **A table-driven test** (`TestRemoveNewLineSuffix`) — a single test function that iterates over a map of cases, running each as a subtest with `t.Run` and `t.Parallel()`.

The table-driven version is the idiomatic Go approach: it keeps the test cases together in one place, runs them as parallel subtests, and scales cleanly as more cases are added.

## Running the tests

```bash
go test -v -run TestRemoveNewLineSuffix
```

To run with race detection (useful for verifying the `t.Parallel()` subtests):

```bash
go test -race -run TestRemoveNewLineSuffix
```

## Related examples in this chapter

- `82-categorizing-tests/` — using `t.Skip` and build tags to categorize tests
- `86-sleeping/` — avoiding `time.Sleep` in tests
- `87-time-api/` — testing time-dependent code
- `88-utility-package/` — test helpers
- `89-benchmark/` — writing benchmarks
- `90-testing-features/` — other testing package features
