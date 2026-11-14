package main

import (
	"fmt"
	"sync"
	"time"
)

func listing1() {
	type Donation struct {
		mu      sync.Mutex
		balance int
	}
	donation := &Donation{}

	// Listener goroutines — polling anti-pattern: unlock, sleep, re-lock.
	// Works, but burns CPU and is slow to react.
	f := func(goal int) {
		donation.mu.Lock()
		for donation.balance < goal {
			donation.mu.Unlock()
			time.Sleep(time.Millisecond * 10)
			donation.mu.Lock()
		}
		fmt.Printf("$%d goal reached\n", donation.balance)
		donation.mu.Unlock()
	}
	go f(2)
	go f(5)

	// Updater goroutine
	go func() {
		for {
			time.Sleep(time.Second)
			donation.mu.Lock()
			donation.balance++
			donation.mu.Unlock()
		}
	}()
}

func listing2() {
	type Donation struct {
		balance int
		ch      chan int
	}

	donation := &Donation{ch: make(chan int)}

	// Listener goroutines
	f := func(goal int) {
		for balance := range donation.ch {
			if balance >= goal {
				fmt.Printf("$%d goal reached\n", balance)
				return
			}
		}
	}
	go f(2)
	go f(5)

	// Updater goroutine
	for {
		time.Sleep(time.Second)
		donation.balance++
		donation.ch <- donation.balance
	}
}

func listing3() {
	type Donation struct {
		cond    *sync.Cond
		balance int
	}

	donation := &Donation{
		cond: sync.NewCond(&sync.Mutex{}),
	}

	// Listener goroutines
	f := func(goal int) {
		donation.cond.L.Lock()
		for donation.balance < goal {
			donation.cond.Wait()
		}
		fmt.Printf("%d$ goal reached\n", donation.balance)
		donation.cond.L.Unlock()
	}
	go f(2)
	go f(5)

	// Updater goroutine
	for {
		time.Sleep(time.Second)
		donation.cond.L.Lock()
		donation.balance++
		donation.cond.L.Unlock()
		donation.cond.Broadcast()
	}
}

func main() {
	_ = listing1 // polling anti-pattern
	_ = listing2 // channel notification
	listing3()   // sync.Cond (recommended)

	// Keep main alive so the goroutines can run and print.
	select {}
}
