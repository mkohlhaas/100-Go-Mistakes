package main

import "fmt"

func listing1() {
	a := [3]int{0, 1, 2}
	for i, v := range a {
		a[2] = 10
		if i == 2 {
			fmt.Println(v) // 2
		}
	}
	fmt.Println(a) // [0 1 10]
}

func listing2() {
	a := [3]int{0, 1, 2}
	for i := range a {
		a[2] = 10
		if i == 2 {
			fmt.Println(a[2]) // 10
		}
	}
	fmt.Println(a) // [0 1 10]
}

func listing3() {
	a := [3]int{0, 1, 2}
	for i, v := range &a {
		a[2] = 10
		if i == 2 {
			fmt.Println(v) // 10
		}
	}
	fmt.Println(a) // [0 1 10]
}

func main() {
	listing1()
	listing2()
	listing3()
}
