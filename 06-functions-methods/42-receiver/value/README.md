# Value receiver (Mistake #42)

This example shows a common Go mistake: mutating state inside a method with a
**value receiver**. Changes made to the receiver only affect a local copy and
are lost when the method returns.

## The code

```go
type customer struct {
	balance float64
}

func (c customer) add(v float64) {
	c.balance += v   // mutates a copy; change is discarded
}

func main() {
	c := customer{balance: 100.}
	c.add(50.)
	fmt.Printf("balance: %.2f\n", c.balance) // balance: 100.00
}
```

`c.add(50.)` is called on `c`, but because `add` has a **value receiver**
`(c customer)`, Go copies `c` into the method. The `balance += v` mutation
applies to that copy only, so `c.balance` stays `100.00` instead of becoming
`150.00`.

## The fix

Use a **pointer receiver** so the method operates on the original struct:

```go
func (c *customer) add(v float64) {
	c.balance += v   // mutates the original; change persists
}
```

See the sibling [`../pointer/`](../pointer/) directory for the corrected
version, which prints `balance: 150.00`.

## When to use each

- Use a **pointer receiver** (`*T`) when the method needs to mutate the
  receiver, or when the struct is large and copying is expensive.
- Use a **value receiver** (`T`) when the method only reads the receiver and
  you want callers to be able to invoke it on both values and pointers
  without explicit `&`.

**Rule of thumb:** keep the receiver type consistent across all methods of a
type — either all value or all pointer receivers.

## Run it

```sh
go run main.go
```

Expected output:

```
balance: 100.00
```