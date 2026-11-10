package main

import "fmt"

func main() {
	fmt.Println("=== handleOperations1 (mistake: nil check) ===")
	handleOperations1("")   // id empty → getOperations returns empty non-nil slice
	handleOperations1("abc") // id non-empty → getOperations returns populated slice

	fmt.Println("\n=== handleOperations2 (correct: len check) ===")
	handleOperations2("")   // id empty → getOperations returns empty non-nil slice
	handleOperations2("abc") // id non-empty → getOperations returns populated slice
}

func handleOperations1(id string) {
	operations := getOperations(id)
	if operations != nil {
		fmt.Printf("  handleOperations1: calling handle with %d operations\n", len(operations))
		handle(operations)
	} else {
		fmt.Println("  handleOperations1: skipping handle (nil slice)")
	}
}

func handleOperations2(id string) {
	operations := getOperations(id)
	if len(operations) != 0 {
		fmt.Printf("  handleOperations2: calling handle with %d operations\n", len(operations))
		handle(operations)
	} else {
		fmt.Println("  handleOperations2: skipping handle (empty slice)")
	}
}

func getOperations(id string) []float32 {
	operations := make([]float32, 0)

	if id == "" {
		return operations // non-nil but empty
	}

	// Add elements to operations
	operations = append(operations, 1.1, 2.2, 3.3)

	return operations
}

func handle(operations []float32) {
	fmt.Printf("    handle: processing %d operations: %v\n", len(operations), operations)
}
