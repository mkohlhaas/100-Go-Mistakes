# Type Embedding and `encoding/json`

Part of [100 Go Mistakes](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them),
Chapter 10 (Standard Library), Mistake 77 — JSON handling.

## What this demonstrates

Embedding `time.Time` (or any type that implements `json.Marshaler`) in a struct
can silently change how `encoding/json` serializes the outer type — and drop
data in the process.

The program compares three ways of pairing an `ID` with a timestamp:

| Listing | Struct | JSON output |
|---------|--------|-------------|
| 1 | `Event1` — embeds `time.Time` anonymously | `"2026-..."` (ID lost!) |
| 2 | `Event2` — `Time time.Time` as a named field | `{"ID":1234,"Time":"2026-..."}` |
| 3 | `Event3` — embeds `time.Time` + custom `MarshalJSON` | `{"ID":1234,"Time":"2026-..."}` |

## The problem: listing 1

```go
type Event1 struct {
    ID int
    time.Time
}
```

`time.Time` implements `json.Marshaler`. Because methods are promoted from an
embedded field, `Event1` also satisfies `json.Marshaler` through the promoted
`time.Time.MarshalJSON`. When you call `json.Marshal(event)`, Go calls that
promoted method, which serializes **only the time** as a JSON string. The `ID`
field is never inspected and is silently dropped:

```text
"2026-11-15T15:21:38.783194668+01:00"   // {"ID":1234} is gone
```

This is easy to miss because the code compiles and runs without error — the
payload simply arrives at the consumer missing a field.

## The fixes

### Listing 2 — don't embed; use a named field (recommended)

```go
type Event2 struct {
    ID   int
    Time time.Time
}
```

`encoding/json` treats `Time` as an ordinary exported field and serializes both
fields normally:

```json
{"ID":1234,"Time":"2026-11-15T15:21:38.783329155+01:00"}
```

This is the clearest, most maintainable approach. Prefer composition (named
fields) over embedding when you want the embedded type's data to appear in the
JSON output.

### Listing 3 — embed, but take over marshaling explicitly

If you must embed `time.Time` (for example to reuse its methods), define
`MarshalJSON` on the outer type and build the output yourself:

```go
type Event3 struct {
    ID int
    time.Time
}

func (e Event3) MarshalJSON() ([]byte, error) {
    return json.Marshal(
        struct {
            ID   int
            Time time.Time
        }{
            ID:   e.ID,
            Time: e.Time,
        },
    )
}
```

Because `Event3` now defines its own `MarshalJSON`, that method shadows the
promoted one from `time.Time`, and both fields are included:

```json
{"ID":1234,"Time":"2026-11-15T15:21:38.783372886+01:00"}
```

## How to run

```bash
go run main.go
```

## Key takeaways

- **Embedding promotes methods.** If an embedded type implements
  `json.Marshaler` (like `time.Time`), the outer struct inherits that behavior,
  and `encoding/json` may use it instead of reflecting over the outer struct's
  own fields.
- **The result can be silent data loss.** Listing 1 drops `ID` with no compile
  or runtime error.
- **Prefer named fields over embedding** for types you want serialized to JSON
  (listing 2).
- If embedding is unavoidable, implement `MarshalJSON`/`UnmarshalJSON` on the
  outer type explicitly so you control exactly what gets serialized (listing 3).
- The same issue applies in reverse on unmarshaling: an embedded
  `json.Unmarshaler` can consume the entire JSON value, leaving the outer
  struct's other fields at their zero values.