package main

import "testing"

var global string

func BenchmarkConcatV1(b *testing.B) {
	var local string
	s := getInput()

	for b.Loop() {
		local = concat1(s)
	}
	global = local
}

func BenchmarkConcatV2(b *testing.B) {
	var local string
	s := getInput()

	for b.Loop() {
		local = concat2(s)
	}
	global = local
}

func BenchmarkConcatV3(b *testing.B) {
	var local string
	s := getInput()

	for b.Loop() {
		local = concat3(s)
	}
	global = local
}

func getInput() []string {
	n := 1_000
	s := make([]string, n)
	for i := range n {
		s[i] = string(make([]byte, 1_000))
	}
	return s
}
