package main

import "testing"

const n = 1_000_000

var global []Bar

func BenchmarkConvert_EmptySlice(b *testing.B) {
	var local []Bar
	foos := make([]Foo, n)

	for b.Loop() {
		local = convertEmptySlice(foos)
	}
	global = local
}

func BenchmarkConvert_GivenCapacity(b *testing.B) {
	var local []Bar
	foos := make([]Foo, n)

	for b.Loop() {
		local = convertGivenCapacity(foos)
	}
	global = local
}

func BenchmarkConvert_GivenLength(b *testing.B) {
	var local []Bar
	foos := make([]Foo, n)

	for b.Loop() {
		local = convertGivenLength(foos)
	}
	global = local
}
