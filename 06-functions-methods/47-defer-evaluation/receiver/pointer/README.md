# defer evaluation with a pointer receiver

This example shows mistake **47** from [*100 Go Mistakes and How to Avoid Them*](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them): the arguments of a deferred call are evaluated when the `defer` statement runs, but with a pointer receiver the call captures the pointer itself — not a copy of the pointed-to data. Mutating the struct after the `defer` therefore changes what the deferred call observes.

## The code

`main.go`

```go
package main

import "fmt"

func main() {
	s := &Struct{id: "foo"}
	defer s.print()
	s.id = "bar"
}

type Struct struct {
	id string
}

func (s *Struct) print() {
	fmt.Println(s.id)
}
```

## Output

```
bar
```

## Explanation

`defer s.print()` evaluates `s.print` at defer time, but `s` is a pointer (`*Struct`). The deferred call keeps referencing that same pointer, so when `main` mutates `s.id` to `"bar"` afterward, the deferred `print()` reads the mutated field and prints `"bar"` — not `"foo"` as one might expect if the struct were captured by value.

## Contrast

Compare with the sibling [`value/`](../value) example: there the receiver is a value (`Struct`), so the defer copies the struct's data at defer time, and printing `"foo"` is the result.

## The fix

Capture the value you want at defer time explicitly:

```go
package main

import "fmt"

func main() {
	s := &Struct{id: "foo"}
	id := s.id
	defer fmt.Println(id)
	s.id = "bar"
}
```

or pass the value as an argument to a closure:

```go
defer func(id string) { fmt.Println(id) }(s.id)
```

## How to run

```sh
go run main.go
```
