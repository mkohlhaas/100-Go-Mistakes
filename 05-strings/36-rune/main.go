package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "hello"
	fmt.Println(len(s)) // 5 — byte count

	s = "€"
	fmt.Println(len(s)) // 3 — byte count, NOT rune count

	// Fixes: count runes instead of bytes
	fmt.Println(utf8.RuneCountInString(s)) // 1

	for _, r := range s {
		fmt.Printf("%c ", r) // €
	}
	fmt.Println()

	// Indexing bytes directly can split a multi-byte rune in two.
	// Convert to []rune to index characters safely.
	runes := []rune(s)
	fmt.Printf("%c\n", runes[0]) // €
}
