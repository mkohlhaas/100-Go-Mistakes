# Pointer Receiver vs Value Receiver with Pointer Fields

This example illustrates a subtle Go gotcha around **method receivers** when a struct contains a **pointer field**.

## The code

```go
type customer struct {
	data *data
}

type data struct {
	balance float64
}

func (c customer) add(operation float64) {
	c.data.balance += operation
}

func main() {
	c := customer{data: &data{
		balance: 100,
	}}
	c.add(50.)
	fmt.Printf("balance: %.2f\n", c.data.balance)
}
```

## What happens

1. `customer` holds a **pointer** to `data`, not `data` by value.
2. `add` uses a **value receiver** `(c customer)`, meaning `c` is a copy of the `customer` struct.
3. Inside `add`, `c.data.balance` is modified.

You might expect the value receiver to protect the original `customer` from mutation — but because `c.data` is a pointer, the copy `c` and the original share the **same** underlying `data`. So `c.data.balance += operation` writes through the pointer and mutates the original balance.

## Output

```
balance: 150.00
```

The balance goes from `100` to `150`, even though `add` has a value receiver.

## The lesson

A value receiver copies the struct, but **does not copy the data pointed to by pointer fields**. If your struct contains pointers, maps, slices, or channels, a value receiver can still mutate shared state.

Use a **pointer receiver** `func (c *customer) add(...)` when the method is intended to mutate the receiver (or when you want to be explicit about it), not because the value receiver fails to mutate pointer fields — it doesn't.

## Run it

```bash
go run main.go
```

## Related

Part of the [100 Go Mistakes](https://www.manning.com/books/100-go-mistakes-and-how-to-avoid-them) series — mistake #42: using the wrong receiver type.
