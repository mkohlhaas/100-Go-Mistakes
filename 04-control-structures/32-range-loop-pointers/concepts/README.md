# Range loop pointers — concepts

Example from *100 Go Mistakes and How to Avoid Them*, mistake #32:
a `range` loop **copies each element into the loop variable**, so taking the
address of the loop variable (`&v`) does not give you a pointer to the
collection element — it gives you a pointer to the loop variable itself, which
is reused on every iteration (pre-Go 1.22). Mutating through that pointer
mutates the same variable over and over, not each element.

## The code

```go
package concepts

type Foo struct{}

type Store struct {
	m map[string]*Foo
}

func (s Store) Put(id string, foo *Foo) {
	s.m[id] = foo
	// ...
}

type LargeStruct struct {
	foo string
	_   [1024]int64
}

func updateMapValue(mapValue map[string]LargeStruct, id string) {
	value := mapValue[id]
	value.foo = "bar"
	mapValue[id] = value
}

func updateMapPointer(mapPointer map[string]*LargeStruct, id string) {
	mapPointer[id].foo = "bar"
}
```

## What the code demonstrates

### Pointers as map values

`Store.m` is a `map[string]*Foo` — the values are **pointers**. Because the
map holds pointers (not copies), callers can mutate the pointed-to `Foo`
through the map, and `Put` can store a pointer even though `Put` has a value
receiver (`s` is copied, but the underlying map — a reference type — is shared).

#### Why a value receiver can still write to the map

`Put` takes `Store` by value, so `s` is a copy of the struct — but `m` is a
map, and maps are reference types: copying the struct copies only the map
header (a pointer to the underlying hash table), not the table itself. So `s`
and the original `store` share the same underlying map.

That's why `s.m[id] = foo` inside `Put` persists on the original `store`, and
why `len(store.m)` reports `2`:

```
Store has 2 entries (pointers stored by value receiver)
```

In short: "pointers stored by value receiver" means that even though `Put`
receives its own copy of the `Store`, the map it points to is shared, so any
entry added through `Put` shows up on the original store. This is the cheap,
mutable-in-place pattern that `updateMapPointer` relies on later in the file.

#### Reference types in Go

Go does not have a formal "reference type" category in the language spec —
**everything in Go is passed by value**. But certain types behave like
references because their value is a header that points to shared underlying
data. Copying the type copies only the header, not the data:

| Type | What the header holds | What's shared |
|------|----------------------|---------------|
| Map | Pointer to the hash table | The underlying entries |
| Slice | Pointer to an array, plus `len` and `cap` | The underlying elements |
| Channel | Pointer to channel state | The channel's buffer / queue |
| Func | Pointer to the function / closure | The function itself |
| Pointer | Address of a variable | The variable at that address |
| Interface | Type pointer + data pointer | The concrete value stored inside |

This is why a map stored in a value-receiver method can still be written to:
`s.m[id] = foo` mutates the shared hash table through the copied map header.
But you **cannot rebind** the caller's variable to a different map — that
requires passing a pointer to the map itself (e.g. `*map[string]int`).

### Map values are copies (the `LargeStruct` problem)

In `updateMapValue`, `mapValue[id]` yields a **copy** of the `LargeStruct`.
Mutating `value.foo` only mutates the copy, so the change is discarded unless
you write it back with `mapValue[id] = value`.

```go
value := mapValue[id]   // copy of the whole LargeStruct (expensive!)
value.foo = "bar"       // mutates the copy only
mapValue[id] = value    // write the copy back — required
```

The `_   [1024]int64` field is what makes `LargeStruct` large:

* `[1024]int64` is a fixed-size Go array of 1024 `int64` values (~8 KiB, since each `int64` is 8 bytes).
* `_` is the blank identifier, so this is an **unnamed field** — you cannot read or write it by name; it simply occupies space in the struct.
* Its only purpose here is **padding**: it inflates `LargeStruct` to ~8 KiB so that storing it as a **map value** makes every copy expensive.

That padding is why `updateMapValue` copies the whole struct when it does `mapValue[id]`, while `updateMapPointer` (which stores `*LargeStruct`) only copies the 8-byte pointer.

### Pointers avoid the copy

In `updateMapPointer`, the map holds `*LargeStruct`, so `mapPointer[id]`
already points at the stored element. Mutating `mapPointer[id].foo` writes
directly to the original — no copy, no write-back needed:

```go
mapPointer[id].foo = "bar" // mutates the original element directly
```

## The mistake

Taking the address of a range loop variable and expecting it to point at each
element in turn:

```go
// ⚠️ BUG (pre-Go 1.22): &v is the same variable every iteration
for _, v := range sliceOfLarge {
    pointers = append(pointers, &v) // all entries alias the same v
    v.foo = "bar"                   // mutates the copy, not the element
}
```

### Go version matters here

**Pre-Go 1.22:** The range loop reuses a **single** loop variable across all
iterations. `&v` is therefore the **same address every time**, so every entry in
`pointers` aliases the same variable. They all end up holding the value of the
**last** iteration:

```
ptrs[0] = 0x..., name="c", val=3   // not "a"!
ptrs[1] = 0x..., name="c", val=3   // all three point to the same variable
ptrs[2] = 0x..., name="c", val=3
```

**Go 1.22+:** Each iteration gets its **own** loop variable, so `&v` differs
per iteration and each pointer addresses a distinct copy:

```
ptrs[0] = 0x..., name="a", val=1   // distinct addresses, distinct copies
ptrs[1] = 0x..., name="b", val=2
ptrs[2] = 0x..., name="c", val=3
```

This program runs on Go 1.27, so it shows the post-1.22 behavior above. But
even with distinct variables, mutating `v` still mutates a **copy** of the
element, not the element itself — that part of the mistake is unchanged.

## The fix

When you need to modify elements inside a `range` loop, avoid copying the
element:

* For **maps**, store pointers (`map[string]*LargeStruct`) and mutate through
  the pointer, as in `updateMapPointer`.
* For **slices**, mutate through the index rather than the copied value:

```go
for i := range sliceOfLarge {
    sliceOfLarge[i].foo = "bar" // mutates the original element
}
```

Or range over a slice of pointers so the loop variable is already a pointer to
the element.

## Run it

```bash
go run main.go
```

Example output:

```
Store has 2 entries (pointers stored by value receiver)
After updateMapValue: "bar" (unchanged — copy was mutated)
After updateMapPointer: "bar" (mutated in place)
After ranging with &it (each ptr addresses a copy):
  ptrs[0] = 0x..., name="a", val=1
  ptrs[1] = 0x..., name="b", val=2
  ptrs[2] = 0x..., name="c", val=3
After mutating by index: items = [{a 10} {b 20} {c 30}]
```

## Reference

- [Go spec — For statements with range](https://go.dev/ref/spec#For_statements)
- [Go Blog — Go 1.22 loop variable semantics](https://go.dev/blog/loopvar-semantics)
- This is example **32** ("range loop pointers") from *100 Go Mistakes and How to Avoid Them*.
