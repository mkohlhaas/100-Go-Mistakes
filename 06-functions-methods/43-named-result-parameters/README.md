# Named result parameters

This example accompanies chapter **43 — Named result parameters** of
*100 Go Mistakes and How to Avoid Them*.

## What are named result parameters?

Go lets you name a function's return values in its signature:

```go
func ReadFull(r io.Reader, buf []byte) (n int, err error)
```

The names `n` and `err` are declared in the function's scope and initialized
to their zero values. A bare `return` statement then returns the current
values of those named results — this is called a **naked return**.

## What this example shows

* `f` — a single named result returned by a naked return.
* `locator` / `loc` — named results used in a method that implements an
  interface (note: the interface itself does not name its results; only the
  implementation does, which is legal).
* `ReadFull` — accumulates bytes read and the last error into named results,
  returning them with a bare `return`, mirroring `io.ReadFull`-style helpers.
* `addNames` — naming a result (`sum`) to document intent at the call site.

## How ReadFull works

`ReadFull` reads into `buf` until it is full or an error occurs:

```go
for len(buf) > 0 && err == nil {
    nr, err = r.Read(buf)
    n += nr
    buf = buf[nr:]
}
```

### What `buf = buf[nr:]` does

It **re-slices** `buf` — it does not copy any data. It creates a new slice
header pointing into the same underlying array, starting `nr` bytes further in:

| Before `buf = buf[nr:]` | After `buf = buf[nr:]` |
|---|---|
| pointer: `&buf[0]` | pointer: `&buf[nr]` |
| length: `L` | length: `L - nr` |
| capacity: `C` | capacity: `C - nr` |

Example: if `buf` is 10 bytes and `r.Read` read 3 bytes, `buf = buf[3:]` makes
the slice start at index 3. The next `r.Read(buf)` writes into `buf[3]` onward
instead of overwriting what was just read. Each iteration shrinks the window by
`nr`, so `len(buf)` eventually reaches 0 and the loop exits.

## When to use them

Named results are most helpful when:

1. The function is long and has multiple return statements — naming the
   results makes each bare `return` self-documenting.
2. The names act as documentation for what each return value means.

## Pitfalls to avoid

### Shadowing

Declaring a local variable with the same name as a named result creates a new
variable that shadows it. The bare `return` still returns the named result
(often its zero value), not the local:

```go
// BUG: returns 0, not 10
func doublePitfall(x int) (result int) {
    {
        result := x * 2 // new local shadows named result
        _ = result
    }
    return // named result result is still 0
}

// FIXED: returns 10
func doubleFixed(x int) (result int) {
    result = x * 2 // assigns to named result, no new variable
    return
}
```

```go
doublePitfall(5) // 0  (want 10)
doubleFixed(5)   // 10
```

The shadowing is most dangerous with `err`: if you write `nr, err := r.Read(buf)`
inside `ReadFull`, the local `err` shadows the named result. The named result
stays `nil`, the loop condition `err == nil` never becomes false, and you can
loop forever — or return a nil error to the caller when the read actually
failed.

### Readability

Naked returns can obscure *what* is being returned, especially in longer
functions. Prefer explicit `return x, y` when the function body is non-trivial.

### API documentation

Named results are not part of a function's type, so they don't affect
assignability or interface satisfaction, but they do show up in godoc —
choose clear names.

## Running

```bash
go run main.go
go vet ./...
go test ./...
```
