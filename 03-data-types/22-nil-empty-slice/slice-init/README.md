# Slice initialization — nil vs. empty slice

This sample demonstrates **Mistake #22 from "100 Go Mistakes": confusing a nil slice with an empty slice.**

## The code

```go
package main

import (
	"fmt"
)

func main() {
	var s []string
	log(1, s)

	s = []string(nil)
	log(2, s)

	s = []string{}
	log(3, s)

	s = make([]string, 0)
	log(4, s)
}

func log(i int, s []string) {
	fmt.Printf("%d: empty=%t\tnil=%t\n", i, len(s) == 0, s == nil)
}
```

## Output

```
1: empty=true	nil=true
2: empty=true	nil=true
3: empty=true	nil=false
4: empty=true	nil=false
```

## What it shows

In Go, **length zero does not mean nil**. All four slices have `len == 0` (so `empty=true`), but they differ in whether they are actually `nil`:

| Line | Initialization | `len == 0` | `s == nil` | Meaning |
|------|----------------|------------|------------|---------|
| 1 | `var s []string` | true | true | nil slice |
| 2 | `s = []string(nil)` | true | true | nil slice |
| 3 | `s = []string{}` | true | false | empty, non-nil slice |
| 4 | `s = make([]string, 0)` | true | false | empty, non-nil slice |

A **nil slice** has no underlying array. An **empty slice** has a (possibly zero-length) non-nil underlying array.

## Why the distinction matters

Even though `len(s) == 0` is true for both, nil and empty slices behave differently:

* **JSON marshaling:** `[]string(nil)` serializes to `null`, while `[]string{}` serializes to `[]`.
* **Explicit nil checks:** `if s == nil` treats them differently, even though ranging over either is safe (iterating a nil slice is a no-op).

## The takeaway

Be deliberate about which one you use:

* Use `var s []string` / `nil` to represent "no value" (absence of a slice).
* Use `[]string{}` / `make([]string, 0)` to represent "an explicitly empty collection".

Don't assume they are interchangeable just because both have length zero.