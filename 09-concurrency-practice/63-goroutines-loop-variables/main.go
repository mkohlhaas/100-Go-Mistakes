package main

import (
	"fmt"
	"sync"
)

func listing1() {
	s := []int{1, 2, 3}
	var wg sync.WaitGroup

	for _, i := range s {
		wg.Go(func() {
			fmt.Print(i)
		})
	}
	wg.Wait()
}

func listing2() {
	s := []int{1, 2, 3}
	var wg sync.WaitGroup

	for _, i := range s {
		val := i
		wg.Go(func() {
			fmt.Print(val)
		})
	}
	wg.Wait()
}

func listing3() {
	s := []int{1, 2, 3}
	var wg sync.WaitGroup

	for _, i := range s {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			fmt.Print(val)
		}(i)
	}
	wg.Wait()
}

func main() {
	fmt.Print("listing1: ")
	listing1()
	fmt.Println()

	fmt.Print("listing2: ")
	listing2()
	fmt.Println()

	fmt.Print("listing3: ")
	listing3()
	fmt.Println()
}
