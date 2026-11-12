package main

import "fmt"

func main() {
	s := "hêllo"

	// Mistake: indexing a string with s[i] yields a BYTE, not a rune.
	// In UTF-8, 'ê' is encoded as two bytes (0xC3 0xAA), so s[1] and s[2]
	// are byte fragments of 'ê', not the character itself.
	fmt.Println("Indexing by byte (the mistake):")
	for i := 0; i < len(s); i++ {
		fmt.Printf("position %d: %c  (byte 0x%x)\n", i, s[i], s[i])
	}
	fmt.Printf("len(s) = %d bytes (there are 5 runes)\n\n", len(s))

	// Correct: ranging over a string decodes UTF-8 and yields each rune
	// along with the byte position where that rune starts.
	fmt.Println("Ranging over the string (correct):")
	for i, r := range s {
		fmt.Printf("position %d: %c\n", i, r)
	}

	// Convert to []rune once, then index by rune position in O(1).
	runes := []rune(s)
	fmt.Println("\nIndexing runes directly (correct):")
	for i := range runes {
		fmt.Printf("rune %d: %c\n", i, runes[i])
	}
}
