package main

import (
	"fmt"
	"time"
)

func listing1() {
	i := 0
	go func() {
		i++
	}()
}

func listing2() {
	i := 0
	go func() {
		i++
	}()
	fmt.Println(i)
}

func listing3() {
	i := 0
	ch := make(chan struct{})
	go func() {
		<-ch
		fmt.Println(i)
	}()
	i++
	ch <- struct{}{}
}

func listing4() {
	i := 0
	ch := make(chan struct{})
	go func() {
		<-ch
		fmt.Println(i)
	}()
	i++
	close(ch)
}

func listing5() {
	i := 0
	ch := make(chan struct{}, 1)
	go func() {
		i = 1
		<-ch
	}()
	ch <- struct{}{}
	fmt.Println(i)
}

func listing6() {
	i := 0
	ch := make(chan struct{})
	go func() {
		i = 1
		<-ch
	}()
	ch <- struct{}{}
	fmt.Println(i)
}

func main() {
	listings := []struct {
		number int
		fn     func()
		// racy listings are run several times to make the
		// nondeterminism visible.
		racy bool
	}{
		{1, listing1, true}, // no output
		{2, listing2, true},
		{3, listing3, false},
		{4, listing4, true},
		{5, listing5, true},
		{6, listing6, false},
	}

	const runs = 5
	for _, listing := range listings {
		fmt.Printf("=== listing %d ===\n", listing.number)
		times := 1
		if listing.racy {
			times = runs
		}
		for i := 0; i < times; i++ {
			listing.fn()
			// Give the goroutine a chance to run so its
			// output (if any) becomes visible before we move on.
			time.Sleep(80 * time.Millisecond)
		}
		fmt.Println()
	}
}
