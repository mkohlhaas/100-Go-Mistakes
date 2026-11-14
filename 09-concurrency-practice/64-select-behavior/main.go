package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=================================")
	fmt.Println("=== buffered channel (cap 10) ===")
	fmt.Println("=================================")

	fmt.Println("\n------------------------------------------------")
	fmt.Println("--- listing1: exit immediately on disconnect ---")
	fmt.Println("------------------------------------------------")
	runListing(10, listing1)

	fmt.Println("\n-----------------------------------------------------")
	fmt.Println("--- listing2: drain messages before disconnecting ---")
	fmt.Println("-----------------------------------------------------")
	runListing(10, listing2)

	fmt.Println("\n--------------------------------------------------------")
	fmt.Println("--- listing3: drain before disconnecting, simplified ---")
	fmt.Println("--------------------------------------------------------")
	runListing(10, listing3)

	fmt.Println("\n==========================")
	fmt.Println("=== unbuffered channel ===")
	fmt.Println("==========================")

	fmt.Println("\n------------------------------------------------")
	fmt.Println("--- listing1: exit immediately on disconnect ---")
	fmt.Println("------------------------------------------------")
	runListing(0, listing1)

	fmt.Println("\n-----------------------------------------------------")
	fmt.Println("--- listing2: drain messages before disconnecting ---")
	fmt.Println("-----------------------------------------------------")
	runListing(0, listing2)

	fmt.Println("\n--------------------------------------------------------")
	fmt.Println("--- listing3: drain before disconnecting, simplified ---")
	fmt.Println("--------------------------------------------------------")
	runListing(0, listing3)
}

func runListing(bufferSize int, fn func(messageCh <-chan int, disconnectCh chan struct{})) {
	messageCh := make(chan int, bufferSize)
	disconnectCh := make(chan struct{})

	go fn(messageCh, disconnectCh)

	for i := range 10 {
		messageCh <- i
	}
	disconnectCh <- struct{}{}
	time.Sleep(10 * time.Millisecond)
}

func listing1(messageCh <-chan int, disconnectCh chan struct{}) {
	for {
		select {
		case v := <-messageCh:
			fmt.Println(v)
		case <-disconnectCh:
			fmt.Println("disconnection, return")
			return
		}
	}
}

func listing2(messageCh <-chan int, disconnectCh chan struct{}) {
	for {
		select {
		case v := <-messageCh:
			fmt.Println(v)
		case <-disconnectCh:
			for {
				select {
				case v := <-messageCh:
					fmt.Println(v)
				default:
					fmt.Println("disconnection, return")
					return
				}
			}
		}
	}
}

func listing3(messageCh <-chan int, disconnectCh chan struct{}) {
	// Phase 1: process messages until disconnect is signaled.
OUTER:
	for {
		select {
		case v := <-messageCh:
			fmt.Println(v)
		case <-disconnectCh:
			break OUTER
		}
	}

	// Phase 2: drain remaining buffered messages before shutting down.
	for {
		select {
		case v := <-messageCh:
			fmt.Println(v)
		default:
			fmt.Println("disconnection, return")
			return
		}
	}
}
