# Utility Function Testing Anti-Pattern

This example demonstrates a common testing mistake in Go: passing `*testing.T` into a utility/helper function.

## The problem

`createCustomer2` takes a `*testing.T` parameter and calls `t.Fatal` directly:

```go
func createCustomer2(t *testing.T, someArg string) Customer {
	customer, err := customerFactory(someArg)
	if err != nil {
		t.Fatal(err) // couples helper to the test framework
	}
	return customer
}
```

This couples the helper to the `testing` package, making it impossible to reuse outside of tests (e.g., in another test package, a `main` program, a benchmark, or a different test framework).

## `t.Helper()`: better error reports, not less coupling

Calling `t.Helper()` at the top of a function tells Go's test framework to treat the function as a test helper. When the helper later calls `t.Fatal`, `t.Error`, `t.FailNow`, etc., the reported failure line number points to the **caller's** line in the test — not to a line inside the helper.

Without `t.Helper()`:

```go
func createCustomer2(t *testing.T, someArg string) Customer {
	customer, err := customerFactory(someArg)
	if err != nil {
		t.Fatal(err) // reported here, inside createCustomer2
	}
	return customer
}
```

With `t.Helper()`:

```go
func createCustomer2(t *testing.T, someArg string) Customer {
	t.Helper() // failures are reported at the call site in the test
	customer, err := customerFactory(someArg)
	if err != nil {
		t.Fatal(err) // reported at the line that called createCustomer2
	}
	return customer
}
```

This makes test output easier to read: when a test fails, `go test` shows the line in `TestCustomer2` that called the helper, so you don't have to hunt through the helper's internals to find what went wrong.

### What `t.Helper()` does not do

`t.Helper()` only affects **where failures are reported**. It does **not**:

- Decouple the helper from `*testing.T`. The helper still imports and depends on the `testing` package.
- Let you reuse the helper outside of tests — not in another package, not in a `main` program, not in a benchmark, and not under a different test framework.
- Turn a test-specific helper into a general-purpose utility.

So `t.Helper()` is worthwhile when the helper is genuinely test-specific and you want nicer failure reports. It is not a substitute for returning errors and letting the caller decide (see [The better approach](#the-better-approach)). Prefer the error-returning style for utilities that might ever be used outside of tests.

## The better approach

Return errors to the caller and let the test decide how to handle them:

```go
func createCustomer1(someArg string) (Customer, error) {
	customer, err := customerFactory(someArg)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}
```

The test then handles the error itself:

```go
func TestCustomer1(t *testing.T) {
	customer, err := createCustomer1("foo")
	if err != nil {
		t.Fatal(err) // the test decides to fail
	}
	_ = customer
}
```

## Helper vs test functions

- **Test functions**: `TestCustomer1`, `TestCustomer2` — start with `Test`, take `*testing.T`, and are run directly by `go test`.
- **Helper functions**: `createCustomer1`, `createCustomer2`, `customerFactory` — called by tests, not run by `go test` directly.

`createCustomer2` is a helper, but it is a *test-specific* helper because it takes `*testing.T`. `createCustomer1` is a plain utility function that works anywhere.

## Test files in the same package

In Go, a `_test.go` file may be in the **same package** as the code under test (internal test):

```go
// main_test.go
package main
```

This is idiomatic and necessary when you need to access unexported identifiers, like `Customer.id` or `customerFactory` in this example.

Go also supports **external tests** in a separate package:

```go
// main_test.go
package main_test
```

External tests can only use the public API, which is useful for verifying the contract consumers depend on. Many Go projects use both: external tests for the public API, and internal tests selectively for unexported logic.

## Testing internals

A strict "unit tests should never test internals" rule is too strong for Go:

- The Go standard library uses internal tests to exercise unexported functions.
- Some internal branches are hard to trigger through public methods alone.
- In a `package main` binary, there may be no meaningful public API beyond `main()`, so testing internals is the only way to test anything.

**Practical guidance:** prefer testing behavior through the public API when you can, but don't treat testing internals as a sin. Keep helpers decoupled from `*testing.T` regardless of where they live.

## Rule of thumb

Utility functions that aren't tests themselves should not accept `*testing.T`. Return errors (or values) and let the caller decide what to do on failure.

## Run the tests

```bash
go test ./...
```
