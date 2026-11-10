# nil vs empty slice — `encoding/json` gotcha

This example (mistake #22 from *100 Go Mistakes*) demonstrates that **a nil slice and an empty slice are not the same value**, and `encoding/json` serializes them differently.

## The code

```go
package main

import (
    "encoding/json"
    "fmt"
)

func main() {
    var s1 []float32                       // nil slice
    customer1 := customer{
        ID:         "foo",
        Operations: s1,
    }
    b, _ := json.Marshal(customer1)
    fmt.Println(string(b))

    s2 := make([]float32, 0)               // empty (non-nil) slice
    customer2 := customer{
        ID:         "bar",
        Operations: s2,
    }
    b, _ = json.Marshal(customer2)
    fmt.Println(string(b))
}

type customer struct {
    ID         string
    Operations []float32
}
```

## Output

```json
{"ID":"foo","Operations":null}
{"ID":"bar","Operations":[]}
```

## Why they differ

| Declaration | Value | JSON output |
|---|---|---|
| `var s1 []float32` | **nil slice** — slice header pointer is `nil` | `null` |
| `make([]float32, 0)` | **empty slice** — backing array allocated, length/capacity are 0 | `[]` |

Both slices have `len() == 0`, but `s1 == nil` while `s2 != nil`. Go's `json.Marshal` treats these as distinct: nil becomes JSON `null`, while an empty slice becomes JSON `[]`.

## Why it matters

If an API client expects `"Operations":[]` but receives `"Operations":null` (or vice versa), it can break deserialization or client logic. The two are not interchangeable at the serialization boundary even though they are "empty" in Go terms.

## The fix

Normalize nil slices to empty slices before marshalling (or pick one convention and document it):

```go
if customer.Operations == nil {
    customer.Operations = []float32{}
}
```

### Key lesson

Nil and empty slices look the same (`len() == 0`) but are different values. Be explicit about which one you mean, especially when crossing serialization boundaries like JSON.