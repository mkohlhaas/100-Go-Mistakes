# Reduce allocations — compiler-optimized map lookup

This example (mistake #96 in *100 Go Mistakes*) shows how the Go compiler
can avoid a heap allocation when converting a `[]byte` to a `string` for a
map lookup — but only when the conversion happens inline inside the map
index expression.

## The code

```go
//go:noinline
func (c *cache) get1(bytes []byte) (v int, contains bool) {
	key := string(bytes)
	v, contains = c.m[key]
	return
}

//go:noinline
func (c *cache) get2(bytes []byte) (v int, contains bool) {
	v, contains = c.m[string(bytes)]
	return
}
```

## The mistake

`get1` converts `bytes` to a `string` and stores it in a named variable
`key` before the map lookup. That intermediate `string` is a value that
could escape to the heap, so in older Go versions (or in more complex call
sites where the compiler can't prove otherwise) every call to `get1`
allocates.

## The fix

`get2` performs the conversion inline inside the map index expression
`c.m[string(bytes)]`. The Go compiler recognizes this pattern and can look
the key up directly from the byte slice without allocating an intermediate
string on the heap — no allocation occurs.

`get2` remains the preferred form: there is never a named intermediate
string that could escape to the heap.

## How to verify

Run escape analysis on the package:

```sh
go build -gcflags="-m" .
```

Sample output:

```
./main.go:28:13: inlining call to fmt.Printf
./main.go:10:7:  c does not escape                   // get1 receiver
./main.go:10:22: bytes does not escape               // get1 arg
./main.go:11:16: string(bytes) does not escape       // get1: no alloc
./main.go:17:7:  c does not escape                   // get2 receiver
./main.go:17:22: bytes does not escape               // get2 arg
./main.go:18:27: string(bytes) does not escape       // get2: no alloc
./main.go:23:7:  &cache{...} does not escape
./main.go:23:31: map[string]int{...} does not escape
./main.go:25:30: [][]byte{...} does not escape
./main.go:25:31: []byte{...} does not escape
./main.go:25:58: []byte{...} does not escape
./main.go:28:13: ... argument does not escape
./main.go:28:60: string(key) escapes to heap         // ESCAPES (fmt formatting)
./main.go:28:61: string(key) escapes to heap         // ESCAPES (fmt formatting)
./main.go:28:67: v1 escapes to heap                  // ESCAPES
./main.go:28:71: ok1 escapes to heap                 // ESCAPES
./main.go:28:76: v2 escapes to heap                  // ESCAPES
./main.go:28:80: ok2 escapes to heap                 // ESCAPES
```

Both functions are marked `//go:noinline`, so the compiler cannot inline
them away. The `string(bytes)` conversions do not escape in either case on
modern Go — the heap allocations come from `string(key)` and the `v`/`ok`
results being passed to `fmt.Printf` in `main`, not from the map lookup.

To observe the historical allocation difference, run the same code on an
older Go version, or make `get1`'s `key` escape (for example, by returning
it or storing it somewhere the compiler can't prove is stack-allocatable).

## Run it

```sh
go run main.go
```

Output:

```
key="hello" get1=(1, true) get2=(1, true)
key="world" get1=(2, true) get2=(2, true)
```
