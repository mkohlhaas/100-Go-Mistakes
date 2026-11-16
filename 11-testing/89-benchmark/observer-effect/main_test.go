package main

import "testing"

const rows = 1000

var res int64

func BenchmarkCalculateSum512_1(b *testing.B) {
	var sum int64
	s := createMatrix512(rows)

	for b.Loop() {
		sum = calculateSum512(s)
	}
	res = sum
}

func BenchmarkCalculateSum513_1(b *testing.B) {
	var sum int64
	s := createMatrix513(rows)

	for b.Loop() {
		sum = calculateSum513(s)
	}
	res = sum
}

func BenchmarkCalculateSum512_2(b *testing.B) {
	var sum int64
	for b.Loop() {
		b.StopTimer()
		s := createMatrix512(rows)
		b.StartTimer()
		sum = calculateSum512(s)
	}
	res = sum
}

func BenchmarkCalculateSum513_2(b *testing.B) {
	var sum int64
	for b.Loop() {
		b.StopTimer()
		s := createMatrix512(rows)
		b.StartTimer()
		sum = calculateSum512(s)
	}
	res = sum
}

func createMatrix512(r int) [][512]int64 {
	return make([][512]int64, r)
}

func createMatrix513(r int) [][513]int64 {
	return make([][513]int64, r)
}
