# iotest

Example from *100 Go Mistakes and How to Avoid Them*, chapter **11 – Testing**, item **88 – Utility package**.

This example demonstrates using Go's standard library `testing/iotest` package to test readers:

- `iotest.TestReader` — verifies that an `io.Reader` behaves correctly by reading known input and checking the output matches expectations.
- `iotest.TimeoutReader` — wraps a reader so reads time out, useful for testing code that handles reader timeouts/errors.

## Files

- `main.go` — example code using `io.ReadAll` and a custom retrying `readAll` helper.
- `main_test.go` — tests exercising `LowerCaseReader`, `foo1`, and `foo2` with `testing/iotest` helpers.

## Run the tests

```bash
go test ./...
```

## Requirements

Go 1.16 or later (`testing/iotest` was added in Go 1.16).
