# monotonic-clock

Example demonstrating a common JSON-handling gotcha in Go: `time.Time` values that carry a monotonic clock reading do not round-trip through `encoding/json` unchanged.

## What it shows

`time.Now()` returns a `Time` that includes both a wall-clock reading **and** a monotonic clock reading (see the monotonic reading (m=...) in the output). When printed, the monotonic reading shows up as an `m=...` suffix (e.g. `2026-11-15 15:05:20.569439127 +0100 CET m=+0.000039442`).

`encoding/json` only serializes the wall-clock portion, so after a marshal/unmarshal round-trip the `m=...` monotonic component is lost.

As a result, two `time.Time` values that appear equal (`event1 == event2`) before serialization can compare as `false` after deserialization.

## The two listings

* `listing1()` — marshals `time.Now()` to JSON and unmarshals it back.  
  Prints `false`, because `event1.Time` still holds the monotonic reading while `event2.Time` does not.

* `listing2()` — calls `t.Truncate(0)` before marshaling, which strips the monotonic reading.  
  Prints `true`, because both sides now compare equal after the round-trip.

`main()` runs both listings in order, printing `false` then `true`.

## Running

```bash
go run main.go
```

Expected output (timestamps will vary):

```
=================
=== Listing 1 ===
=================
Original event:     {2026-11-15 15:16:51.542268461 +0100 CET m=+0.000052599}
Unmarshalled event: {2026-11-15 15:16:51.542268461 +0100 CET}
false

=================
=== Listing 2 ===
=================
Original event:     {2026-11-15 15:16:51.542462213 +0100 CET}
Unmarshalled event: {2026-11-15 15:16:51.542462213 +0100 CET}
true
```

Notice listing1's `Original event` carries `m=...` (monotonic reading present), but its `Unmarshalled event` does not — `encoding/json` dropped it — hence `false`. listing2 strips the reading up front with `Truncate(0)`, so neither line carries `m=...` — hence `true`.

## The fix

Strip the monotonic clock reading before serializing or comparing `time.Time` values:

```go
t := time.Now().Truncate(0)
```

Or, when you need to preserve monotonic measurements, avoid comparing wall-clock-only values across a JSON boundary — carry the monotonic information separately instead.

## Reference

Part of the [100 Go Mistakes](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them) series — standard library / JSON handling.
