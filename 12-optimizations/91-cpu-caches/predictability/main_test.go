package main

import "testing"

var global int64

const n = 1_000_000

func BenchmarkLinkedList(b *testing.B) {
	var local int64
	for b.Loop() {
		b.StopTimer()
		nodes := make([]node, n)
		for i := range n - 1 {
			nodes[i].next = &nodes[i+1]
		}
		b.StartTimer()
		local = linkedList(&nodes[0])
	}
	global = local
}

func BenchmarkSum2(b *testing.B) {
	var local int64
	for b.Loop() {
		b.StopTimer()
		s := make([]int64, n)
		b.StartTimer()
		local = sum2(s)
	}
	global = local
}
