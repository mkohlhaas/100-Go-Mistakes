package main

import (
	"fmt"
	"runtime"
)

func consumeMessages(n int, getType func([]byte) []byte) {
	for range n {
		msg := receiveMessage()
		// Do something with msg
		storeMessageType(getType(msg))
	}
}

func getMessageType(msg []byte) []byte {
	return msg[:5]
}

func getMessageTypeWithCopy(msg []byte) []byte {
	msgType := make([]byte, 5)
	copy(msgType, msg)
	return msgType
}

func receiveMessage() []byte {
	return make([]byte, 1_000_000)
}

var stored [][]byte

func storeMessageType(msgType []byte) {
	stored = append(stored, msgType)
}

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d KB\n", m.Alloc/1024)
}

func main() {
	const iterations = 100

	// Leaky path: reslicing shares the 1 MB backing array.
	// Each retained 5-byte slice pins its entire 1 MB buffer, so GC cannot reclaim it.
	fmt.Println("\n--- Leaky: getMessageType (msg[:5]) ---")
	fmt.Print("Alloc before: ")
	printAlloc()
	consumeMessages(iterations, getMessageType)
	fmt.Printf("Alloc after %d messages: ", iterations)
	printAlloc()
	runtime.GC()
	fmt.Printf("Alloc after GC (still leaked): ")
	printAlloc()

	// Fixed path: copying breaks the shared backing array.
	// Only the small 5-byte copies are retained; the 1 MB buffers are freed.
	fmt.Println("--- Fixed: getMessageTypeWithCopy ---")
	fmt.Print("Alloc before: ")
	printAlloc()
	consumeMessages(iterations, getMessageTypeWithCopy)
	fmt.Printf("Alloc after %d messages: ", iterations)
	printAlloc()
	runtime.GC()
	fmt.Printf("Alloc after GC: ")
	printAlloc()
}
