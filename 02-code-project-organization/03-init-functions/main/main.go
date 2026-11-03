package main

import (
	"fmt"

	"github.com/teivah/100-go-mistakes/02-code-project-organization/03-init-functions/redis"
)

// Two init() functions in the same package run in source-file order
// (init 1, then init 2) before main() executes.
//
// MISTAKE: this ordering is implicit and fragile. It depends on file names and
// source order, and can silently change when files are renamed, reordered, or
// new init() functions are added.
func init() {
	fmt.Println("init 1")
}

func init() {
	fmt.Println("init 2")
}

func main() {
	err := redis.Store("foo", "bar")
	_ = err
}
