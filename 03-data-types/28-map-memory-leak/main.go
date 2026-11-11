package main

import (
	"fmt"
	"runtime"
)

func main() {
	// Init
	n := 1_000_000
	m := make(map[int][128]byte)
	printAlloc() // 0

	// Add n elements
	for i := range n {
		m[i] = randBytes()
	}
	printAlloc() // 294MB

	// Remove n elements
	for i := range n {
		delete(m, i)
	}

	// End
	runtime.GC()
	printAlloc() // ⚠️ 287MB
	runtime.KeepAlive(m)
}

func randBytes() [128]byte {
	return [128]byte{}
}

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d MB\n", m.Alloc/1024/1024)
}
