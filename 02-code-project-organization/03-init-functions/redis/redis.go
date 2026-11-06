// Package redis
package redis

import "fmt"

// init runs automatically when the redis package is imported, printing
// "redis" before any caller code runs.
//
// MISTAKE: importing redis for Store() triggers this side effect even though
// the caller only needs Store. Package initialization should not perform
// observable work that the importer did not explicitly request.
func init() {
	fmt.Println("redis")
}

// Store is a no-op stub. In real code this would store a key/value pair.
func Store(key, value string) error {
	return nil
}
