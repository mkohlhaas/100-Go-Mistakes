# `strings.Trim` vs `strings.TrimSuffix` / `strings.TrimPrefix`

This example illustrates a common mistake in Go: confusing `strings.Trim` / `TrimLeft` / `TrimRight` with `strings.TrimSuffix` / `strings.TrimPrefix`.

## The mistake

`Trim`, `TrimLeft`, and `TrimRight` treat their second argument as a **cutset** (a set of characters), not as a substring/suffix/prefix. They strip *any* character contained in the cutset from the given side(s), repeating until a character not in the set is reached.

`TrimSuffix` and `TrimPrefix`, by contrast, remove the argument as an **exact substring** at most once.

## Example

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.TrimRight("123oxo", "xo"))  // 123
	fmt.Println(strings.TrimSuffix("123oxo", "xo")) // 123o
	fmt.Println(strings.TrimLeft("oxo123", "ox"))   // 123
	fmt.Println(strings.TrimPrefix("oxo123", "ox")) // o123
	fmt.Println(strings.Trim("oxo123oxo", "ox"))    // 123
}
```

**Output:**

```
123
123o
123
o123
123
```

## Why the results differ

| Call | What happens | Result |
|------|--------------|--------|
| `TrimRight("123oxo", "xo")` | Strips every trailing `x` or `o`: `o`, `x`, `o` are all removed until `3` | `123` |
| `TrimSuffix("123oxo", "xo")` | Removes the exact suffix `xo` once | `123o` |
| `TrimLeft("oxo123", "ox")` | Strips every leading `o` or `x`: `o`, `x`, `o` are all removed until `1` | `123` |
| `TrimPrefix("oxo123", "ox")` | Removes the exact prefix `ox` once | `o123` |
| `Trim("oxo123oxo", "ox")` | Strips leading and trailing `o`/`x` characters from both sides | `123` |

## Rule of thumb

- Use `TrimSuffix` / `TrimPrefix` when you want to remove a specific **substring** from the end / beginning of a string.
- Use `TrimRight` / `TrimLeft` / `Trim` when you want to remove **any of several characters** (a cutset) from the end / beginning / both sides.

## Running

```sh
go run main.go
```

## Reference

- [`strings.Trim`](https://pkg.go.dev/strings#Trim)
- [`strings.TrimLeft`](https://pkg.go.dev/strings#TrimLeft)
- [`strings.TrimRight`](https://pkg.go.dev/strings#TrimRight)
- [`strings.TrimPrefix`](https://pkg.go.dev/strings#TrimPrefix)
- [`strings.TrimSuffix`](https://pkg.go.dev/strings#TrimSuffix)
