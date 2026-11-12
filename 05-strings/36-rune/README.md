# Mistake 36 — Confusing bytes and runes

Go strings are read-only byte slices, not sequences of characters. A single
Unicode rune (code point) can be encoded in 1 to 4 bytes, so `len(s)` returns
the **byte count**, not the number of runes.

## What this example shows

```go
s := "hello"
fmt.Println(len(s)) // 5 — 5 bytes, 5 runes

s = "€"
fmt.Println(len(s)) // 3 — 3 bytes, but only 1 rune
```

Output:

```
5
3
```

`"€"` is 3 bytes long in UTF-8, even though it is a single character.

## The fix

To count runes, iterate over the string with `range` (which yields runes) or
use `utf8.RuneCountInString`. To index characters safely, convert to `[]rune`
first — indexing a string directly gives you bytes, which can split a
multi-byte rune in two:

```go
import "unicode/utf8"

s := "€"
fmt.Println(utf8.RuneCountInString(s)) // 1

for _, r := range s {
    fmt.Printf("%c ", r) // €
}

runes := []rune(s)
fmt.Printf("%c\n", runes[0]) // €
```

## Run it

```bash
go run main.go
```

## See also

- [Mistake 36 in the book](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them)
- [`unicode/utf8`](https://pkg.go.dev/unicode/utf8) package