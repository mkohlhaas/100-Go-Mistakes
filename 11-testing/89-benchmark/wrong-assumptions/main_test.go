package main

import (
	"sync/atomic"
	"testing"
)

func BenchmarkAtomicStoreInt32(b *testing.B) {
	var v atomic.Int32
	for b.Loop() {
		v.Store(1)
	}
}

func BenchmarkAtomicStoreInt64(b *testing.B) {
	var v atomic.Int64
	for b.Loop() {
		v.Store(1)
	}
}
