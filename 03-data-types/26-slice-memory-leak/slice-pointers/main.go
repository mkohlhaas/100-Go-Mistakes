package main

import (
	"fmt"
	"runtime"
)

type Foo struct {
	v []byte
}

func main() {
	scenario1Reslice()
	scenario2Copy()
	scenario3MarkNil()
}

// scenario1Reslice keeps only the first two elements via a reslice.
// The reslice shares the original backing array, so all 1000 byte slices
// stay reachable and nothing is freed until 'two' is dropped.
func scenario1Reslice() {
	printSeparator("Scenario 1: keepFirstTwoElementsOnly (reslice)")
	fmt.Println("Allocates 1000 × 1 MB, then keeps only foos[:2]. Because the reslice shares the backing array, all 1000 byte slices stay reachable and nothing is freed.")

	foos := makeFooSlice()
	printAlloc() // ~1024333 KB (~1 GB allocated)

	two := keepFirstTwoElementsOnly(foos)
	printAlloc() // ~1024336 KB (reslice shares backing array; still ~1 GB)

	runtime.GC()
	printAlloc() // ~1024336 KB (leaked: 'two' keeps the whole backing array alive)

	runtime.KeepAlive(two) // prevent 'two' from being optimized away

	runtime.GC()
	printAlloc() // ~336 KB (released once 'two' and 'foos' are dead per liveness analysis)
}

// scenario2Copy copies only the first two elements into a new slice.
// After GC, the original backing array becomes unreachable and is collected.
func scenario2Copy() {
	printSeparator("Scenario 2: keepFirstTwoElementsOnlyCopy (copy)")
	fmt.Println("Allocates 1000 × 1 MB, then copies only the first two elements into a new slice. After GC, the original backing array can be collected.")

	foos := makeFooSlice()
	printAlloc() // ~1024339 KB (~1 GB allocated)

	two := keepFirstTwoElementsOnlyCopy(foos)
	printAlloc() // ~1024341 KB (copy adds only ~2 KB)

	runtime.GC()
	printAlloc() // ~2389 KB (original backing array freed; only ~2 KB remains)

	runtime.KeepAlive(two)

	runtime.GC()
	printAlloc() // ~341 KB
}

// scenario3MarkNil nils out elements 2–999 before reslicing, breaking the
// references so the unused byte slices can be collected after GC.
func scenario3MarkNil() {
	printSeparator("Scenario 3: keepFirstTwoElementsOnlyMarkNil (mark nil)")
	fmt.Println("Allocates 1000 × 1 MB, then nils out elements 2–999 before reslicing. After GC, the unused elements can be collected.")

	foos := makeFooSlice()
	printAlloc() // ~1024339 KB (~1 GB allocated)

	two := keepFirstTwoElementsOnlyMarkNil(foos)
	printAlloc() // ~1024341 KB (marking nil adds no allocation)

	runtime.GC()
	printAlloc() // ~2389 KB (unused elements freed)

	runtime.KeepAlive(two)

	runtime.GC()
	printAlloc() // ~341 KB
}

func printSeparator(title string) {
	fmt.Println()
	fmt.Println("########################################################################")
	fmt.Println("#", title)
	fmt.Println("########################################################################")
	fmt.Println()
}

func makeFooSlice() []Foo {
	foos := make([]Foo, 1_000)
	for i := range foos {
		foos[i] = Foo{v: make([]byte, 1024*1024)}
	}
	return foos
}

func keepFirstTwoElementsOnly(foos []Foo) []Foo {
	return foos[:2]
}

func keepFirstTwoElementsOnlyCopy(foos []Foo) []Foo {
	res := make([]Foo, 2)
	copy(res, foos)
	return res
}

func keepFirstTwoElementsOnlyMarkNil(foos []Foo) []Foo {
	for i := 2; i < len(foos); i++ {
		foos[i].v = nil
	}
	return foos[:2]
}

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("alloc: %d KB\n", m.Alloc/1024)
}
