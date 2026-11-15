package main

import (
	"fmt"
	"time"
)

func main() {
	listing1()
	_ = listing2
}

func listing1() {
	ticker := time.NewTicker(1000)
	for range ticker.C {
		fmt.Println("tick#1")
	}
}

func listing2() {
	ticker := time.NewTicker(time.Microsecond)
	for range ticker.C {
		fmt.Println("tick#2")
	}
}
