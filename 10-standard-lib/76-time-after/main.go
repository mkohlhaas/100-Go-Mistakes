package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// warnAfter is how long a consumer waits for a message before warning.
// It defaults to time.Hour (the mistake as presented in the book) but
// main shortens it so the demo actually triggers the warning quickly.
var warnAfter = time.Hour

func consumer1(ch <-chan Event) {
	for {
		select {
		case event := <-ch:
			handle(1, event)
		case <-time.After(warnAfter):
			// In every iteration of the for loop a new timer is instantiated
			// and never stopped — that's the timer leak this mistake demonstrates.
			log.Println("consumer1 warning: no messages received")
		}
	}
}

func consumer2(ch <-chan Event) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), warnAfter)
		select {
		case event := <-ch:
			cancel()
			handle(2, event)
		case <-ctx.Done():
			log.Println("consumer2 warning: no messages received")
		}
	}
}

func consumer3(ch <-chan Event) {
	timerDuration := warnAfter
	timer := time.NewTimer(timerDuration)

	for {
		timer.Reset(timerDuration)
		select {
		case event := <-ch:
			handle(3, event)
		case <-timer.C:
			log.Println("consumer3 warning: no messages received")
		}
	}
}

type Event struct{}

func handle(i int, _ Event) {
	fmt.Printf("Handling an Event: consumer %d\n", i)
}

// producer sends n events on ch, spaced by delay, then stops sending.
// The channel is deliberately left open so that consumers can keep
// blocking on it until the demo shuts them down.
func producer(ch chan<- Event, n int, delay time.Duration) {
	for i := range n {
		select {
		case ch <- Event{}:
			log.Printf("produced event %d", i+1)
		case <-time.After(delay):
			return
		}
	}
}

func main() {
	// Shorten the timeout for the demo so the "no messages" warning
	// actually fires within a few seconds instead of after an hour.
	warnAfter = 2 * time.Second

	for _, c := range []struct {
		name     string
		consumer func(<-chan Event)
	}{
		{"consumer1 (time.After)", consumer1},
		{"consumer2 (context.WithTimeout)", consumer2},
		{"consumer3 (time.NewTimer + Reset)", consumer3},
	} {
		go func() {
			ch := make(chan Event)
			go producer(ch, 2, 300*time.Millisecond)
			c.consumer(ch)
		}()
	}

	time.Sleep(5 * time.Second)
	log.Println("demo done")
}
