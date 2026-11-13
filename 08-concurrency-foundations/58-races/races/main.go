package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func listing1() int {
	i := 0

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		i++
	}()
	go func() {
		defer wg.Done()
		i++
	}()
	wg.Wait()

	return i
}

func listing2() int64 {
	var i atomic.Int64

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		i.Add(1)
	}()
	go func() {
		defer wg.Done()
		i.Add(1)
	}()
	wg.Wait()

	return i.Load()
}

func listing3() int {
	i := 0
	mutex := sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		mutex.Lock()
		i++
		mutex.Unlock()
	}()
	go func() {
		defer wg.Done()
		mutex.Lock()
		i++
		mutex.Unlock()
	}()
	wg.Wait()

	return i
}

func listing4() int {
	i := 0
	ch := make(chan int)

	go func() {
		ch <- 1
	}()
	go func() {
		ch <- 1
	}()

	i += <-ch
	i += <-ch

	return i
}

func listing5() int {
	i := 0
	mutex := sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		i = 1
	}()
	go func() {
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		i = 2
	}()
	wg.Wait()

	return i
}

func main() {
	const runs = 5

	fmt.Println("=== listing 1: unsynchronized i++ (data race) ===")
	for n := range runs {
		fmt.Printf("  run %d: i = %d (want 2, result is nondeterministic)\n", n+1, listing1())
	}

	fmt.Printf("=== listing 2: sync/atomic ===\n  i = %d (want 2)\n", listing2())
	fmt.Printf("=== listing 3: sync.Mutex ===\n  i = %d (want 2)\n", listing3())
	fmt.Printf("=== listing 4: unbuffered channel ===\n  i = %d (want 2)\n", listing4())
	fmt.Printf("=== listing 5: sync.Mutex, last write wins ===\n  i = %d (1 or 2)\n", listing5())
}
