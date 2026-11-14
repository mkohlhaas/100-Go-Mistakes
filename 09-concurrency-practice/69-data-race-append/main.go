package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== listing1: appending to a shared slice ===")
	listing1()
	time.Sleep(200 * time.Millisecond)

	fmt.Println("=== listing2: appending to private copies ===")
	listing2()
	time.Sleep(200 * time.Millisecond)

	fmt.Println("done")
}

func listing1() {
	s := make([]int, 1)

	go func() {
		s1 := append(s, 1)
		fmt.Println(s1)
	}()

	go func() {
		s2 := append(s, 1)
		fmt.Println(s2)
	}()
}

func listing2() {
	s := make([]int, 0, 1)

	go func() {
		sCopy := make([]int, len(s), cap(s))
		copy(sCopy, s)

		s1 := append(sCopy, 1)
		fmt.Println(s1)
	}()

	go func() {
		sCopy := make([]int, len(s), cap(s))
		copy(sCopy, s)

		s2 := append(sCopy, 1)
		fmt.Println(s2)
	}()
}
