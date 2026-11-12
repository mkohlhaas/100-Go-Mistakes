package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	listing1()
	fmt.Println()

	listing2()
	fmt.Println()

	runListing3()
	runListing4()
}

func listing1() {
	for i := range 5 {
		fmt.Printf("%d ", i) // 0 1 2 3 4

		switch i {
		default:
		case 2:
			break
		}
	}
}

func listing2() {
loop:
	for i := range 5 {
		fmt.Printf("%d ", i) // 0 1 2

		switch i {
		default:
		case 2:
			break loop
		}
	}
}

func listing3(ctx context.Context, ch <-chan int) {
	for {
		select {
		case <-ch:
			// Do something
		case <-ctx.Done():
			// Output: listing3 still spinning after ctx canceled — break exits the select, not the for loop
			break
		}
	}
}

func listing4(ctx context.Context, ch <-chan int) {
loop:
	for {
		select {
		case <-ch:
			// Do something
		case <-ctx.Done():
			// Output: listing4 returned — break loop exits the for loop
			break loop
		}
	}
}

func runListing3() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		listing3(ctx, nil)
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("listing3 returned")
	case <-time.After(150 * time.Millisecond):
		cancel()
		fmt.Println("listing3 still spinning after ctx canceled — break exits the select, not the for loop")
	}
}

func runListing4() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		listing4(ctx, nil)
		close(done)
	}()

	cancel()
	<-done
	fmt.Println("listing4 returned — break loop exits the for loop")
}
