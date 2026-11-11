# Mistake #31: Range loop argument evaluation

This example demonstrates a subtle Go gotcha: the expression in a `for ... range` loop is evaluated **once**, before the loop starts. Reassigning the loop's target variable inside the loop body does not change what is being ranged over.

## The code

```go
package main

import "fmt"

func main() {
	ch1 := make(chan int, 3)
	go func() {
		ch1 <- 0
		ch1 <- 1
		ch1 <- 2
		close(ch1)
	}()

	ch2 := make(chan int, 3)
	go func() {
		ch2 <- 10
		ch2 <- 11
		ch2 <- 12
		close(ch2)
	}()

	ch := ch1
	for v := range ch {
		fmt.Println(v)
		ch = ch2   // <- has no effect on the iteration
	}
}
```

## What happens

1. `ch := ch1` — `ch` points to `ch1`.
2. `for v := range ch` — the channel expression `ch` is evaluated **once**, so the loop ranges over `ch1` for its entire lifetime.
3. Inside the loop, `ch = ch2` reassigns the local variable, but the `range` loop is still iterating over `ch1`. The assignment is ignored by the loop.
4. Once `ch1` is drained and closed, the loop ends. `ch2` is never read.

**Linter output:** `ch = ch2` is flagged as an ineffectual assignment, and the new value of `ch` is never used — precisely because `range` captured `ch1` before the loop began.

**Output:**

```
0
1
2
```

Values `10`, `11`, `12` are never printed.

## The fix

Range over each channel explicitly, or use a `select` loop if you want to read from multiple channels:

```go
for v := range ch1 {
    fmt.Println(v)
}
for v := range ch2 {
    fmt.Println(v)
}
```

Or, to merge two channels:

```go
for {
    select {
    case v, ok := <-ch1:
        if !ok {
            ch1 = nil
        } else {
            fmt.Println(v)
        }
    case v, ok := <-ch2:
        if !ok {
            ch2 = nil
        } else {
            fmt.Println(v)
        }
    }
    if ch1 == nil && ch2 == nil {
        return
    }
}
```

## How to run

```bash
go run main.go
```

## Reference

This is example **31** ("range loop argument evaluation") from *100 Go Mistakes and How to Avoid Them*.
