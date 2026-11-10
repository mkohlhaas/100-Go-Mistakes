package main

import "fmt"

func main() {
	var n float32 = 1.0001
	fmt.Println(n * n) // 1.0002

	var a float64
	positiveInf := 1 / a
	negativeInf := -1 / a
	nan := a / a
	fmt.Println(positiveInf, negativeInf, nan) // +Inf -Inf NaN

	fmt.Println(f1(100), f2(100)) // 10100.00999999993 10100.01
}

func f1(n int) float64 {
	result := 10_000.
	for range n {
		result += 1.0001
	}
	return result
}

func f2(n int) float64 {
	result := 0.
	for range n {
		result += 1.0001
	}
	return result + 10_000.
}
