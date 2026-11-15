package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Time time.Time
}

func listing1() error {
	t := time.Now()
	event1 := Event{
		Time: t,
	}

	fmt.Println("Original event:    ", event1)

	b, err := json.Marshal(event1)
	// fmt.Println("Marshalled event:  ", string(b))
	if err != nil {
		return err
	}

	var event2 Event
	err = json.Unmarshal(b, &event2)
	fmt.Println("Unmarshalled event:", event2)
	if err != nil {
		return err
	}

	fmt.Println(event1 == event2) // false
	return nil
}

func listing2() error {
	t := time.Now()
	event1 := Event{
		Time: t.Truncate(0), // ✅ the fix
	}

	fmt.Println("Original event:    ", event1)

	b, err := json.Marshal(event1)
	// fmt.Println("Marshalled event:  ", string(b))
	if err != nil {
		return err
	}

	var event2 Event
	err = json.Unmarshal(b, &event2)
	fmt.Println("Unmarshalled event:", event2)
	if err != nil {
		return err
	}

	fmt.Println(event1 == event2) // true
	return nil
}

func main() {
	fmt.Println("=================")
	fmt.Println("=== Listing 1 ===")
	fmt.Println("=================")

	if err := listing1(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println()

	fmt.Println("=================")
	fmt.Println("=== Listing 2 ===")
	fmt.Println("=================")

	if err := listing2(); err != nil {
		fmt.Println(err)
		return
	}
}
