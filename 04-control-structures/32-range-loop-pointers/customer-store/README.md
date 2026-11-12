# Customer Store — Range Loop Pointer Mistake

A small Go example demonstrating mistake **#32** from *100 Go Mistakes*:
**storing pointers to range-loop variables in a map**.

## The problem

```go
func (s *Store) storeCustomers(customers []Customer) {
	for _, customer := range customers {
		fmt.Printf("%p\n", &customer)
		s.m[customer.ID] = &customer
	}
}
```

`customer` is a single variable reused across every iteration of the loop. Taking `&customer` always yields the **same address**, so every entry in `s.m` ends up pointing at the last customer's value instead of each individual customer:

```
key=1, value=&main.Customer{ID:"3", Balance:0}
key=2, value=&main.Customer{ID:"3", Balance:0}
key=3, value=&main.Customer{ID:"3", Balance:0}
```

This is the output on **Go ≤1.21**. On **Go ≥1.22** (including the Go 1.27 used here), each range-iteration variable is distinct, so `&customer` is a different address every iteration and this particular bug no longer occurs. That's why running `go run main.go` today shows correct output even for the "buggy" version — the `%p` addresses printed inside the loop confirm each iteration has its own variable:

```
=== storeCustomers (buggy) ===
0x159fbda72048
0x159fbda72060
0x159fbda72078
key=1, value=&main.Customer{ID:"1", Balance:10}
key=2, value=&main.Customer{ID:"2", Balance:-10}
key=3, value=&main.Customer{ID:"3", Balance:0}
=== storeCustomers2 (copy before &) ===
key=3, value=&main.Customer{ID:"3", Balance:0}
key=1, value=&main.Customer{ID:"1", Balance:10}
key=2, value=&main.Customer{ID:"2", Balance:-10}
=== storeCustomers3 (index into slice) ===
key=1, value=&main.Customer{ID:"1", Balance:10}
key=2, value=&main.Customer{ID:"2", Balance:-10}
key=3, value=&main.Customer{ID:"3", Balance:0}
```

The addresses differ per iteration on modern Go, so no aliasing happens. To see the bug reproduce, test against Go ≤1.21: all three map entries then point at customer 3.

## The fixes

**Copy the value before taking its address:**

```go
func (s *Store) storeCustomers2(customers []Customer) {
	for _, customer := range customers {
		current := customer
		s.m[current.ID] = &current
	}
}
```

**Index into the slice directly (preferred):**

```go
func (s *Store) storeCustomers3(customers []Customer) {
	for i := range customers {
		customer := &customers[i]
		s.m[customer.ID] = customer
	}
}
```

`storeCustomers3` avoids an extra copy entirely by taking the address of each slice element directly.

**Note on `current := customer`:** On **Go ≥1.22**, this copy is superfluous — each range iteration already has its own variable, so `&customer` is unique per iteration. `storeCustomers2` could drop the `current := customer` line and still be correct on modern Go. The copy was only required on Go ≤1.21.

`storeCustomers3` remains the preferred fix on all versions because it avoids the copy entirely and points directly into the slice's backing array.

## Run it

```bash
go run main.go
```

`main` calls all three versions so you can compare their output side by side.

## Lesson

Never store `&loopVar` from a `for ... range` loop unless you've first copied the variable into a new local, or you index into the underlying slice/array. This is one of the most common Go mistakes, especially when building maps or slices of pointers inside a loop.
