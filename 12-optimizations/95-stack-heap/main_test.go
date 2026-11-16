package main

import "testing"

var (
	globalValue int
	globalPtr   *int
)

func BenchmarkSumValue(b *testing.B) {
	b.ReportAllocs()
	var local int
	for i := 0; b.Loop(); i++ {
		local = sumValue(i, i)
	}
	globalValue = local
}

func BenchmarkSumPtr(b *testing.B) {
	b.ReportAllocs()
	var local *int
	for i := 0; b.Loop(); i++ {
		local = sumPtr(i, i)
	}
	globalValue = *local
}
