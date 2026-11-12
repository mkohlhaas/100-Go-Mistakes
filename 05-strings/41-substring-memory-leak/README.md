# 41 - Substring memory leak

Example for "substring memory leak" from *100 Go Mistakes and How to Avoid Them*.

## The mistake

In Go, slicing a string does **not** copy its data:

```go
uuid := log[:36]
```

The resulting substring shares the underlying byte array of the original
string. As long as `uuid` is kept alive (e.g. stored in a struct, cached, or
passed around), the **entire** original `log` string stays in memory too —
even if you only needed 36 bytes from a multi-megabyte log line. That is a
memory leak.

## The fixes

`main()` extracts a 36-character UUID from each log line three ways:

| Approach | Code | Verdict |
|----------|------|---------|
| `handleLog1` | `uuid := log[:36]` | ❌ Leaks — keeps the whole log string alive |
| `handleLog2` | `uuid := string([]byte(log[:36]))` | ✅ Copies — `[]byte`→`string` allocates a new backing array |
| `handleLog3` | `uuid := strings.Clone(log[:36])` | ✅ Copies — idiomatic, dedicated function (Go 1.17+) |

Prefer `strings.Clone` for clarity; it says exactly what you intend.

The valid logs are processed by all three handlers, so each UUID appears three
times in the stored list. The `"short"` line fails the length check in every
handler and prints an error instead.

## Byte slicing on runes

`main()` also shows that slicing by bytes cuts multi-byte characters:

```go
s1 = "Hêllo, World!"
s2 = string([]rune(s1)[:5]) // "Hêllo", safe
```

Use rune conversion when slicing text that may contain non-ASCII characters.

## Run it

```sh
go run .
```

Output:

```
Hello
Hêllo
handleLog1: log is not correctly formatted
handleLog2: log is not correctly formatted
handleLog3: log is not correctly formatted
stored UUIDs: [00000000-0000-0000-0000-000000000001 00000000-0000-0000-0000-000000000001 00000000-0000-0000-0000-000000000001 00000000-0000-0000-0000-000000000002 00000000-0000-0000-0000-000000000002 00000000-0000-0000-0000-000000000002 00000000-0000-0000-0000-000000000003 00000000-0000-0000-0000-000000000003 00000000-0000-0000-0000-000000000003]
```

Requires Go 1.17 or later (`strings.Clone`).
