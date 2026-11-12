# defer evaluation with a value receiver

This example demonstrates a common Go gotcha: **deferred function calls are evaluated when the `defer` statement runs, not when the deferred call executes.**

```go
package main

import "fmt"

func main() {
	s := Struct{id: "foo"}
	defer s.print()  // s is copied here (id == "foo")
	s.id = "bar"
}

type Struct struct {
	id string
}

func (s Struct) print() {
	fmt.Println(s.id)
}
```

## Output

```
foo
```

## Why

Because `print` has a **value receiver**, `defer s.print()` copies `s` at the moment the `defer` statement executes. The copy is then invoked when `main` returns, so the later mutation `s.id = "bar"` does not affect it — the deferred call still prints `"foo"`.

To capture the updated value, use a **pointer receiver** or defer a closure that reads `s` at execution time:

```go
func (s *Struct) print() { fmt.Println(s.id) }

// or
defer func() { fmt.Println(s.id) }()
```

## Run

```bash
go run main.go
```
