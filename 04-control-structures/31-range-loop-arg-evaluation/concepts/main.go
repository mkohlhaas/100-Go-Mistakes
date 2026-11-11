package main

import "fmt"

func main() {
	s1 := []int{0, 1, 2}
	for range s1 {
		s1 = append(s1, 10)
	}
	fmt.Println(s1) // [0 1 2 10 10 10]

	s2 := []int{0, 1, 2}
	for i := 0; i < len(s2); i++ {
		s2 = append(s2, 10) // ⚠️ signal: killed (s2 grows infinitely)
	}
	fmt.Println(s2)
}
