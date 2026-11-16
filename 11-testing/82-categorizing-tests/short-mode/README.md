# Categorizing Tests with `testing.Short()`

This example shows how to categorize tests by **runtime cost** using `testing.Short()`, so long-running tests are skipped when you run `go test -short` and only run when you explicitly opt in to the full suite.

## Files

```
main_test.go   # one fast test, one long-running test gated behind testing.Short()
```

## The example

`main_test.go` has two tests. `TestFast` always runs. `TestLongRunning` simulates
expensive work (here a `time.Sleep`) and checks `testing.Short()` first, skipping
itself when short mode is on:

```go
package main

import (
	"testing"
	"time"
)

func TestFast(t *testing.T) {
	t.Log("fast test body")
}

func TestLongRunning(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test")
	}
	t.Log("running long-running test body")
	time.Sleep(250 * time.Millisecond)
	t.Log("long-running test body finished")
}
```

## Running the tests (you can try this now)

### Short mode — fast feedback

```bash
go test -short -v ./...
```

Output:

```
=== RUN   TestFast
    main_test.go:10: fast test body
--- PASS: TestFast (0.00s)
=== RUN   TestLongRunning
    main_test.go:18: skipping long-running test
--- SKIP: TestLongRunning (0.00s)
PASS
ok  	100-go-mistakes/.../short-mode	0.003s
```

`testing.Short()` returns `true`, so `TestLongRunning` calls `t.Skip` and is
reported as skipped. The whole suite finishes in about **3 ms** — the long-running
test never executes.

### Full mode — run everything

```bash
go test -v ./...
```

Output:

```
=== RUN   TestFast
    main_test.go:10: fast test body
--- PASS: TestFast (0.00s)
=== RUN   TestLongRunning
    main_test.go:20: running long-running test body
    main_test.go:22: long-running test body finished
--- PASS: TestLongRunning (0.25s)
PASS
ok  	100-go-mistakes/.../short-mode	0.253s
```

`testing.Short()` returns `false`, so `TestLongRunning` runs its body (including
the 250 ms sleep) and the suite takes about **250 ms**.

### Why `t.Skip` instead of an early `return`?

Using `t.Skip` (rather than just returning) makes the skip visible in verbose
output and in the test summary, so you can see the test was considered and
deliberately skipped — not silently ignored or forgotten.

### Is `-short` predefined? Is it always there?

Yes. `-short` is built into Go's standard `testing` package — it is always
available, `false` by default, and becomes `true` when you pass it:

```bash
go test -short ./...
```

You never define or register the flag yourself; your test just reads it with
`testing.Short()`. You can see all the built-in test flags with `go help testflag`.

## CI recommendation

Run the short suite on every commit/PR so fast feedback stays fast:

```bash
go test -short ./...
```

Run the full suite separately (e.g., nightly or in a dedicated CI job) without
`-short`:

```bash
go test ./...
```

## Why this matters

- Keeps the fast test suite fast — no accidental slow tests blocking quick feedback.
- Lets developers opt in to long-running tests locally when they want to (`go test` without `-short`).
- Makes test runtime cost explicit and visible in CI output via skips.

## See also

The sibling directory `build-tags/` shows a different way to categorize tests —
by **kind** (unit vs. integration) using Go build tags rather than by **runtime
cost** (`-short`).