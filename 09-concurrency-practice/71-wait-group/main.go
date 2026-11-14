package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func listing1() {
	wg := sync.WaitGroup{}
	var v uint64

	for range 3 {
		go func() {
			wg.Add(1)
			atomic.AddUint64(&v, 1)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(v)
}

func listing2() {
	wg := sync.WaitGroup{}
	var v uint64

	wg.Add(3)
	for range 3 {
		go func() {
			atomic.AddUint64(&v, 1)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(v)
}

func listing3() {
	wg := sync.WaitGroup{}
	var v uint64

	for range 3 {
		wg.Go(func() {
			atomic.AddUint64(&v, 1)
		})
	}

	wg.Wait()
	fmt.Println(v)
}

func main() {
	fmt.Println("listing1:")
	listing1() // 0

	fmt.Println("\nlisting2:")
	listing2() // 3

	fmt.Println("\nlisting3:")
	listing3() // 3
}
