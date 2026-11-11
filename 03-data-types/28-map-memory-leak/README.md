# Map Memory Leak

Demonstrates that deleting all entries from a Go map does not free the map's backing memory, even after garbage collection.

## What happens

`main.go` allocates a `map[int][128]byte`, inserts 1,000,000 elements (~294 MB), then deletes every element and runs `runtime.GC()`.

Despite all entries being removed, the process still holds ~287 MB. Go maps never shrink their internal bucket array — deleted entries are marked empty, but the memory stays allocated for the lifetime of the map.

## Output

```
0 MB
294 MB
287 MB
```

## The fix

If you need the memory back, replace the map with a fresh one instead of clearing it in place:

```go
m = make(map[int][128]byte)
```

The old map becomes unreachable and its memory can be reclaimed by the GC (or reused by the runtime).

## Run it

```bash
go run main.go
```

Need Go 1.13+ for the `1_000_000` digit separator.

## Related

This is mistake #28 from [100-go-mistakes](https://github.com/100-go-mistakes/100-go-mistakes): "Assuming maps free memory when entries are deleted."
