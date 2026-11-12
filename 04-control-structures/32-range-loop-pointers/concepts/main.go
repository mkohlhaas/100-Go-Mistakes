// Package concepts demonstrates mistake #32: range loop pointers.
package main

import "fmt"

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

func main() {
	// --- Pointers as map values (Store) ---
	store := Store{m: make(map[string]*Foo)}
	store.Put("a", &Foo{})
	store.Put("b", &Foo{})
	fmt.Printf("Store has %d entries (pointers stored by value receiver)\n", len(store.m)) // Store has 2 entries (pointers stored by value receiver)

	// --- Map values are copies (updateMapValue) ---
	valueMap := map[string]LargeStruct{
		"x": {foo: "baz"},
	}
	updateMapValue(valueMap, "x")
	fmt.Printf("After updateMapValue: %q\n", valueMap["x"].foo) // After updateMapValue: "bar"

	// --- Pointers avoid the copy (updateMapPointer) ---
	pointerMap := map[string]*LargeStruct{
		"y": {foo: "baz"},
	}
	updateMapPointer(pointerMap, "y")
	fmt.Printf("After updateMapPointer: %q (mutated in place)\n", pointerMap["y"].foo) // After updateMapPointer: "bar" (mutated in place)

	// --- Range loop pointer mistake ---
	type Item struct {
		name string
		val  int
	}
	items := []Item{{name: "a", val: 1}, {name: "b", val: 2}, {name: "c", val: 3}}

	var ptrs []*Item
	for _, item := range items {
		ptrs = append(ptrs, &item) // &item points to a copy of the element, not the element itself
	}
	fmt.Println("After ranging with &item (each ptr addresses a copy):")
	for i, p := range ptrs {
		fmt.Printf("  ptrs[%d] = %p, name=%q, val=%d\n", i, p, p.name, p.val)
	}
	// Go 1.22+ behavior.
	// After ranging with &item (each ptr addresses a copy):
	//   ptrs[0] = 0x1cde48e12000, name="a", val=1
	//   ptrs[1] = 0x1cde48e12018, name="b", val=2
	//   ptrs[2] = 0x1cde48e12030, name="c", val=3

	// --- Fix: mutate/index by pointer to the element ---
	for i := range items {
		items[i].val *= 10
	}
	fmt.Println("After mutating by index: items =", items) // After mutating by index: items = [{a 10} {b 20} {c 30}]
}
