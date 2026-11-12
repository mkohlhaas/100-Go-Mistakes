# Mistake 37 — Iterating over strings the wrong way

## The mistake

A Go string is a read-only sequence of **bytes**. Indexing into it with
`s[i]` returns the byte at position `i`, not the character (rune). That
works for ASCII, but breaks for any UTF-8 text: characters like `ê` are
encoded as multiple bytes, so `s[i]` can hand you a byte fragment and
`len(s)` reports bytes, not runes.

```go
s := "hêllo"
fmt.Println(s[1]) // prints a byte of 'ê', not 'ê'
fmt.Println(len(s)) // 6, not 5
```

## The fix

- Range over the string: `for i, r := range s`. Go decodes the UTF-8
  and yields each rune together with the byte position where it starts.
- Convert to `[]rune(s)` when you need rune-based indexing.

## Running the example

```sh
go run main.go
```

## Output

```
Indexing by byte (the mistake):
position 0: h   (byte 0x68)
position 1: Ã   (byte 0xc3)
position 2: ª   (byte 0xaa)
position 3: l   (byte 0x6c)
position 4: l   (byte 0x6c)
position 5: o   (byte 0x6f)
len(s) = 6 bytes (there are 5 runes)
```

> Note: `for i := range s` over a string does **not** visit every byte offset — it visits the start byte of each rune (0, 1, 3, 4, 5 here), which is why position 2 never appears as a rune start. That's exactly why indexing strings by byte offset is misleading.

Ranging over the string (correct):
position 0: h
position 1: ê
position 3: l
position 4: l
position 5: o

Indexing runes directly (correct):
rune 0: h
rune 1: ê
rune 2: l
rune 3: l
rune 4: o
```

## Key takeaway

Never index a string by byte offset to read characters. Decode runes
with `range`, or convert to `[]rune` first.
