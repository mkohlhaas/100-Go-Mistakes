package main

import "testing"

const n = 1_000_000

var global map[int]struct{}

func BenchmarkMapWithoutSize(b *testing.B) {
	var local map[int]struct{}
	for b.Loop() {
		m := make(map[int]struct{})
		for j := range n {
			m[j] = struct{}{}
		}
	}
	global = local
}

func BenchmarkMapWithSize(b *testing.B) {
	var local map[int]struct{}
	for b.Loop() {
		m := make(map[int]struct{}, n)
		for j := range n {
			m[j] = struct{}{}
		}
	}
	global = local
}
