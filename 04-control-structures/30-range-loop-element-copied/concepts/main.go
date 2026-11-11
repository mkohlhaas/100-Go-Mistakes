package main

import (
	"fmt"
	"strings"
)

// ## The mistake
// range copies each slice element into the loop variable by value, so
// mutating the loop variable does not change the original slice element.

func main() {
	s := []string{"a", "b", "c"}
	for i, v := range s {
		fmt.Printf("index=%d, value=%s\n", i, v)
	}

	for i := range s {
		fmt.Printf("index=%d\n", i)
	}

	for _, v := range s {
		fmt.Printf("value=%s\n", v)
	}

	// ⚠️ mutating the range variable does not change the slice
	for _, v := range s {
		v = strings.ToUpper(v)         // mutates only the copy
		fmt.Printf("local: v=%s\n", v) // uppercase locally...
	}
	fmt.Println(s) // [a b c] — unchanged! ...but the slice itself is not modified

	// Fix: mutate through the index
	for i := range s {
		s[i] = strings.ToUpper(s[i])
	}
	fmt.Println(s) // [A B C]
}
