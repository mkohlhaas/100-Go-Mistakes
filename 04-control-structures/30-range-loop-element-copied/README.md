# Range loop element copied

Part of **"100 Go Mistakes"** — Chapter 4: Control Structures, item 30.

## The mistake

In Go, a `range` loop over a slice (or array) **copies each element into the loop variable** by value. Mutating the loop variable mutates only the copy, not the original slice element:

```go
s := []string{"a", "b", "c"}
for _, v := range s {
    v = strings.ToUpper(v) // mutates only the copy
}
fmt.Println(s) // [a b c] — unchanged!
```

This catches people coming from languages where the loop variable is a reference to the collection element.

## What's in this folder

| Directory | What it demonstrates |
|-----------|----------------------|
| `concepts/` | Basic `range` mechanics — index-only, value-only, and index+value forms |
| `value-copy/` | The classic bug: incrementing `balance` on a copied `account` struct has no effect on the original slice |

## Examples

### `concepts/main.go` — how `range` assigns variables

```go
for i, v := range s {
    fmt.Printf("index=%d, value=%s\n", i, v)
}

for i := range s {       // i is the index; v is discarded
    fmt.Printf("value=%d\n", i)
}

for _, v := range s {    // i is discarded; v is the element copy
    fmt.Printf("value=%s\n", v)
}
```

### `value-copy/main.go` — mutation doesn't propagate

```go
accounts := createAccounts() // [{100} {200} {300}]

for _, a := range accounts {
    a.balance += 1000 // mutates the copy only
}
fmt.Println(accounts) // [{100} {200} {300}] — unchanged

// Correct: mutate through the index
for i := range accounts {
    accounts[i].balance += 1000
}
fmt.Println(accounts) // [{1100} {1200} {1300}]

// Pointers avoid the copy, but watch out for aliasing
accountsPtr := createAccountsPtr()
for _, a := range accountsPtr {
    a.balance += 1000 // a is a pointer; mutation propagates
}
```

## The fix

When you need to modify slice elements inside a `range` loop, mutate through the index rather than the copied value:

```go
for i := range accounts {
    accounts[i].balance += 1000
}
```

Or use a slice of pointers so the loop variable points at the original element (be mindful of shared/aliased state in that case).

## When the copy matters

For small types (int, string) the copy is cheap and usually harmless — unless you assign to the variable expecting the slice to change. For large structs, the copy can be expensive, which is another reason to prefer index-based mutation or pointers when working with big elements inside a hot loop.

## Why no `&` in `createAccountsPtr()`?

```go
func createAccountsPtr() []*account {
	return []*account{
		{balance: 100.},  // Go implicitly takes the address: &account{balance: 100.}
		{balance: 200.},
		{balance: 300.},
	}
}
```

You don't need to write `&account{...}` explicitly. When a composite literal is used where a pointer type is expected (here, slice elements are `*account`), Go **automatically takes its address**. `{balance: 100.}` is an `account` value literal, and Go treats it as `&account{balance: 100.}`.

Writing `&account{balance: 100.}` would be identical; the shorthand is just a convenience.

## Run it

```bash
cd concepts   && go run main.go
cd value-copy && go run main.go
```

## Learn more

- [Effective Go — For (range)](https://go.dev/doc/effective_go#for)
- [Go Blog — Go 1.22 loop variable semantics](https://go.dev/blog/loopvar-semantics) (loop variable reuse, a related but separate change in Go 1.22+)
