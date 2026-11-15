# JSON handling with `map[string]any` — Mistake #77

This directory contains example code from chapter "10 - Standard Library", item **77** of
[100 Go Mistakes and How to Avoid Them](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them).

## The mistake

```go
var m map[string]any
err := json.Unmarshal(b, &m)
```

Unmarshaling JSON into `map[string]any` throws away type information:

- JSON numbers are decoded as `float64`, even when they represent integers (`int`, `uint`, `int64`, …).
- There is no way to tell whether `1` came from JSON `1` or `1.0`.
- Refactoring becomes fragile: changing a JSON field's type silently changes the Go type at runtime.

## The fix

Prefer a concrete **struct** so types are preserved and validated:

```go
type Message struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Active bool  `json:"active"`
}

var msg Message
err := json.Unmarshal(b, &msg)
```

If you genuinely don't know the shape up front, use `json.Decoder` with `UseNumber()` to keep numbers as `json.Number` instead of `float64`:

```go
dec := json.NewDecoder(bytes.NewReader(b))
dec.UseNumber()
var m map[string]any
if err := dec.Decode(&m); err != nil {
    return err
}
// m["id"] is now json.Number, convert explicitly when needed
```

## Running the example

From this directory:

```bash
go run main.go
```

Expected output:

```
m = map[string]interface {}{"id":42, "name":"widget", "price":9.99}
m["id"] type: float64, value: 42
BUG: m["id"].(int) failed — JSON numbers are float64, not int
m["id"] == 42 is false (surprising, but correct: float64 vs int)
message.ID type: int, value: 42
message.ID == 42 is true
message = {ID:42 Name:widget Price:9.99}
```

The example builds against the module at the repo root (`go 1.27`):

```bash
go build ./...
```

## Related

- Parent directory: `../` (other JSON-handling examples: `monotonic-clock`, `type-embedding`)
- Book repo root: `../../../..`
