# Map iteration

Example from "100 Go Mistakes" — mistake 33: iterating over a map while modifying it.

## What it shows

Go maps are unordered, and their contents may change while you range over them. Adding entries to a map during iteration is allowed, but entries added during iteration may or may not be produced by the range loop — there is no guarantee.

- **listing1** modifies the map (`m[10+k] = true`) while ranging over it. Because the map grows during iteration, newly added keys may or may not be visited, so the output is non-deterministic.
- **listing2** avoids the problem by iterating over `m` and writing changes into a separate copy (`m2`, built with `maps.Copy`). The output is deterministic.

## Run it

```sh
go run main.go
```

## Fix

Never add or delete entries in the map you are currently ranging over. If you need to update a map based on its own contents, iterate over the original and apply changes to a copy — `maps.Copy` (Go 1.21+) makes the copy one line — or collect the changes first and apply them after the loop finishes.
