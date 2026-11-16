# Reducing allocations with `sync.Pool`

This example shows how to reuse short-lived buffers with `sync.Pool` to reduce heap allocations (from "100 Go Mistakes", mistake #96).

## What it does

`write` borrows a pre-allocated 1 KB byte buffer from the pool, uses it to build a response, then resets and returns it to the pool:

```go
var pool = sync.Pool{
    New: func() any {
        b := make([]byte, 1024)
        return &b
    },
}

func write(w io.Writer) {
    bPtr := pool.Get().(*[]byte)
    defer func() {
        *bPtr = (*bPtr)[:0]
        pool.Put(bPtr)
    }()

    b := *bPtr
    getResponse(b)
    _, _ = w.Write(b)
}
```

What the reset-and-return does:

```go
*bPtr = (*bPtr)[:0]   // reset the slice
pool.Put(bPtr)        // return it to the pool
```

- `bPtr` is a `*[]byte` pulled from the pool.
- `(*bPtr)[:0]` re-slices the underlying array to length `0` while keeping its **capacity**. The slice looks empty to the next caller, but the allocated array is preserved — no new heap allocation is needed next time.
- `*bPtr = ...` writes that emptied slice back through the pointer.
- `pool.Put(bPtr)` hands the pointer back to `sync.Pool` so another goroutine can `Get()` it.

Resetting the length clears stale data from the previous response; keeping the capacity avoids re-allocating.

Without the pool, every call to `write` would allocate a new 1 KB slice on the heap. With the pool, the same slice is reused across calls.

## Run it

```bash
go run main.go
```

## Caveats when using `sync.Pool`

- **Objects may be garbage collected.** `pool.Get()` can return a freshly created object (via `New`) even if nothing was put back, so never assume a pooled object is still valid across goroutines.
- **Never keep a reference to a pooled object after returning it.** The buffer is reset and reused, so retaining the pointer means you'll see data from another caller.
- **Best for short-lived, frequently allocated objects** (buffers, context objects, temp structs). Not a general-purpose object cache — the runtime may discard pooled items at any time.
