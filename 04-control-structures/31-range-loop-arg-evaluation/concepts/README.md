# Range loop argument evaluation — concepts

Example from *100 Go Mistakes and How to Avoid Them*, mistake #31:
the expression in a `for ... range` loop is evaluated **once**, before the
loop starts. Reassigning the slice you are ranging over inside the loop
body does not change what the loop iterates over.

## The code

```go
package main

func main() {
	s1 := []int{0, 1, 2}
	for range s1 {
		s1 = append(s1, 10)
	}

	s2 := []int{0, 1, 2}
	for i := 0; i < len(s2); i++ {
		s2 = append(s2, 10)
	}
}
```

## What happens

**`s1` — `range` evaluates its argument once**

`for range s1` evaluates `s1` exactly once, before iteration begins, so the
loop walks the original slice (length 3) and runs exactly **3 times**, even
though `append` grows `s1` on every iteration:

```
s1 after loop: [0 1 2 10 10 10]
```

The appended elements are never visited, because the range loop is still
iterating over the snapshot taken before the loop started.

**`s2` — a C-style `for` re-evaluates `len(s2)` every iteration**

`for i := 0; i < len(s2); i++ { s2 = append(s2, 10) }` re-reads `len(s2)`
each time, and `append` makes the slice longer every iteration, so the
condition never becomes false. This is an **infinite loop** (it runs until
memory is exhausted).

## The mistake

```go
for range s1 {
    s1 = append(s1, 10) // ⚠️ does not extend the iteration
}
```

Assuming that mutating (or reassigning) the range target inside the body
affects what the loop iterates over. It does not — `range` captured the
slice before the first iteration.

## The fix

If you need to keep appending while iterating, iterate by index over a
snapshot, or collect the results separately:

```go
// iterate over a fixed snapshot
for i := 0; i < len(s1); i++ {
    _ = s1[i]
}

// or build a new slice instead of appending to the one being ranged over
added := make([]int, 0, len(s1))
for range s1 {
    added = append(added, 10)
}
s1 = append(s1, added...)
```

Or simply don't modify the range target inside the loop.

## Run it

```bash
go run main.go
```

> Note: the `s2` loop in `main.go` is infinite. Run the `s1` portion in
> isolation if you want to observe finite output.

## Reference

- [Go spec — For statements with range](https://go.dev/ref/spec#For_statements)
- This is example **31** ("range loop argument evaluation") from *100 Go Mistakes and How to Avoid Them*.
