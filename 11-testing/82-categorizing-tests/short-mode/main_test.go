package main

import (
	"testing"
	"time"
)

// TestFast is a quick unit-style test that always runs, even in short mode.
func TestFast(t *testing.T) {
	t.Log("fast test body")
}

// TestLongRunning simulates an expensive test (e.g. hitting a slow service,
// crunching a large dataset, or retrying with backoff). It skips itself when
// testing.Short() is true so `go test -short` stays fast.
func TestLongRunning(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test")
	}
	t.Log("running long-running test body")
	time.Sleep(250 * time.Millisecond)
	t.Log("long-running test body finished")
}
