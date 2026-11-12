# Pointer Receiver

This example demonstrates how a **pointer receiver** lets a method modify the original struct in place.

## Code

```go
type customer struct {
	balance float64
}

func (c *customer) add(operation float64) {
	c.balance += operation
}

func main() {
	c := customer{balance: 100.0}
	c.add(50.0)
	fmt.Printf("balance: %.2f\n", c.balance) // balance: 150.00
}
```

## What's happening

- `c *customer` means `add` receives a **pointer** to `c`, not a copy.
- Because it operates on the original value, `c.balance` is updated permanently — `main` sees `150.00`.
- If `add` used a **value receiver** `(c customer)`, it would modify a copy and `main` would still print `100.00` (see the `value/` sibling example).

## Why use a pointer receiver

- To mutate the receiver's fields and have the change visible to the caller.
- To avoid copying large structs on every method call.
- Go automatically takes the address when you call a pointer-receiver method on an addressable value (`c.add(...)` works even though `c` is not a `*customer`).

## Run it

```bash
go run main.go
```

Expected output:

```
balance: 150.00
```

## Related examples

- [`../value`](../value) — value receiver: changes are local to a copy.
- [`../struct-with-pointer`](../struct-with-pointer) — value receiver on a struct holding a pointer field: the pointed-to data can still be mutated.
