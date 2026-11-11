package main

import (
	"fmt"
	"reflect"
)

type customer1 struct {
	id string
}

type customer2 struct {
	id         string
	operations []float64
}

func main() {
	cust1_1 := customer1{id: "x"}
	cust1_2 := customer1{id: "x"}
	fmt.Println(cust1_1 == cust1_2) // true

	cust2_1 := customer2{id: "x", operations: []float64{1.}}
	cust2_2 := customer2{id: "x", operations: []float64{1.}}
	// Doesn't compile
	// fmt.Println(cust2_1 == cust2_2) // ⚠️ invalid operation: cust2_1 == cust2_2 (struct containing []float64 cannot be compared) [UndefinedOp]
	_ = cust2_1
	_ = cust2_2

	var a any = 3
	var b any = 3
	fmt.Println(a == b) // true

	// var cust3_1 any = customer2{id: "x", operations: []float64{1.}}
	// var cust3_2 any = customer2{id: "x", operations: []float64{1.}}
	// fmt.Println(cust3_1 == cust3_2) // ⚠️ panic: runtime error: comparing uncomparable type main.customer2

	cust4_1 := customer2{id: "x", operations: []float64{1.}}
	cust4_2 := customer2{id: "x", operations: []float64{1.}}
	fmt.Println(reflect.DeepEqual(cust4_1, cust4_2)) // true
	fmt.Println(cust4_1.equal(cust4_2))              // true

	// different operations → false
	cust4_3 := customer2{id: "x", operations: []float64{2.}}
	fmt.Println(cust4_1.equal(cust4_3)) // false
}

func (a customer2) equal(b customer2) bool {
	if a.id != b.id {
		return false
	}
	if len(a.operations) != len(b.operations) {
		return false
	}
	for i := 0; i < len(a.operations); i++ {
		if a.operations[i] != b.operations[i] {
			return false
		}
	}
	return true
}
