package main

import (
	"fmt"
	"math"
)

func Inc32(counter int32) int32 {
	if counter == math.MaxInt32 {
		panic("int32 overflow")
	}
	return counter + 1
}

func IncInt(counter int) int {
	if counter == math.MaxInt {
		panic("int overflow")
	}
	return counter + 1
}

func IncUint(counter uint) uint {
	if counter == math.MaxUint {
		panic("uint overflow")
	}
	return counter + 1
}

func AddInt(a, b int) int {
	if a > math.MaxInt-b {
		panic("int overflow")
	}

	return a + b
}

func MultiplyInt(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}

	result := a * b
	if a == 1 || b == 1 {
		return result
	}
	if a == math.MinInt || b == math.MinInt {
		panic("integer overflow")
	}
	if result/b != a {
		panic("integer overflow")
	}
	return result
}

// try runs f and recovers from any panic, printing the result or the panic message.
func try(name string, f func()) {
	fmt.Printf("--- %s ---\n", name)
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  panic: %v\n", r)
			}
		}()
		f()
	}()
	fmt.Println()
}

func main() {
	// --- Inc32 ---
	try("Inc32(100) - safe", func() {
		fmt.Printf("Inc32(100) = %d\n", Inc32(100))
	})
	try("Inc32(math.MaxInt32) - overflow", func() {
		fmt.Printf("Inc32(math.MaxInt32) = %d\n", Inc32(math.MaxInt32))
	})

	// --- IncInt ---
	try("IncInt(100) - safe", func() {
		fmt.Printf("IncInt(100) = %d\n", IncInt(100))
	})
	try("IncInt(math.MaxInt) - overflow", func() {
		fmt.Printf("IncInt(math.MaxInt) = %d\n", IncInt(math.MaxInt))
	})

	// --- IncUint ---
	try("IncUint(100) - safe", func() {
		fmt.Printf("IncUint(100) = %d\n", IncUint(100))
	})
	try("IncUint(math.MaxUint) - overflow", func() {
		fmt.Printf("IncUint(math.MaxUint) = %d\n", IncUint(math.MaxUint))
	})

	// --- AddInt ---
	try("AddInt(10, 20) - safe", func() {
		fmt.Printf("AddInt(10, 20) = %d\n", AddInt(10, 20))
	})
	try("AddInt(math.MaxInt, 1) - overflow", func() {
		fmt.Printf("AddInt(math.MaxInt, 1) = %d\n", AddInt(math.MaxInt, 1))
	})

	// --- MultiplyInt ---
	try("MultiplyInt(6, 7) - safe", func() {
		fmt.Printf("MultiplyInt(6, 7) = %d\n", MultiplyInt(6, 7))
	})
	try("MultiplyInt(math.MaxInt, 2) - overflow", func() {
		fmt.Printf("MultiplyInt(math.MaxInt, 2) = %d\n", MultiplyInt(math.MaxInt, 2))
	})
	try("MultiplyInt(math.MinInt, -1) - overflow", func() {
		fmt.Printf("MultiplyInt(math.MinInt, -1) = %d\n", MultiplyInt(math.MinInt, -1))
	})
}
