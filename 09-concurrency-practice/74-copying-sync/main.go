package main

import (
	"sync"
	"time"
)

func main() {
	counter1 := NewCounter1()

	go func() {
		counter1.Increment1("foo")
	}()
	go func() {
		counter1.Increment1("bar")
	}()

	time.Sleep(10 * time.Millisecond)

	counter2 := NewCounter2()

	go func() {
		counter2.Increment2("foo")
	}()
	go func() {
		counter2.Increment2("bar")
	}()
}

// -------- //
// Counter1 //
// -------- //

type Counter1 struct {
	mu       sync.Mutex // ⚠️ bad
	counters map[string]int
}

func NewCounter1() Counter1 {
	return Counter1{counters: map[string]int{}}
}

// Increment1 ⚠️ passes lock by value
func (c Counter1) Increment1(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[name]++
}

// -------- //
// Counter2 //
// -------- //

type Counter2 struct {
	mu       *sync.Mutex // 👍 good
	counters map[string]int
}

func NewCounter2() Counter2 {
	return Counter2{
		mu:       &sync.Mutex{},
		counters: map[string]int{},
	}
}

// Increment2 - pass Counter2 by reference (pointer receiver)
func (c *Counter2) Increment2(name string) {
	// Same code
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[name]++
}
