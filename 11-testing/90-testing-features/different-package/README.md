# Testing from a Different Package (External Test Package)

This example demonstrates **Mistake #90: Testing a package from within itself instead of from an external test package** from *100 Go Mistakes and How to Avoid Them*.

When a test file is in the same package as the code under test (a so-called "internal test", `package counter`), it can access unexported identifiers. That sounds convenient, but it couples the tests to the package's internals and lets the test bypass the public API. Writing tests in a different package (`package counter_test`, an "external test") forces you to exercise only the exported API, which produces tests that behave like real consumers of the package and refactor more safely.

## Files

```
main.go           # library under test: package counter
main_test.go      # tests: package counter_test (external test package)
```

## The example

`main.go` defines a thread-safe counter:

```go
// Package counter
package counter

import "sync/atomic"

var count uint64

func Inc() uint64 {
	atomic.AddUint64(&count, 1)
	return count
}
```

`main_test.go` tests it from an **external** package:

```go
package counter_test

import (
	"testing"

	counter "100-go-mistakes/11-testing/90-testing-features/different-package"
)

func TestCount(t *testing.T) {
	if counter.Inc() != 1 {
		t.Errorf("expected 1")
	}
}
```

Notice:

- The test declares `package counter_test`, not `package counter`. The `_test` suffix tells Go this is an external test package.
- It imports the `counter` package and calls only the exported `Inc` function. It cannot (and does not) touch the unexported `count` variable.

## Running the tests

From this directory:

```bash
go test ./...
```

Output:

```
ok  	100-go-mistakes/11-testing/90-testing-features/different-package	0.002s
```

With verbose output to see the test run:

```bash
go test -v .
```

Output:

```
=== RUN   TestCount
--- PASS: TestCount (0.0s)
PASS
ok  	100-go-mistakes/11-testing/90-testing-features/different-package	0.002s
```

## Why this matters

- **Tests use the public API**, exactly like any other consumer of the package. If you later rename or remove unexported internals, the tests keep passing as long as the exported contract is stable.
- **No white-box coupling**: internal tests can reach into unexported state, which makes tests brittle and tempts you to assert on implementation details instead of behavior.
- **Avoids import cycles**: an internal test file is part of the package, so it cannot import anything that imports the package back; an external test package is a normal importer and never creates a cycle.

## Guideline

Prefer `package <name>_test` for your test files. Only use an internal test package (`package <name>`) when you genuinely need to test an unexported symbol that has no exported equivalent, and treat that as the exception.
