package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.TrimRight("123oxo", "xo"))  // 123
	fmt.Println(strings.TrimSuffix("123oxo", "xo")) // 123o
	fmt.Println(strings.TrimLeft("oxo123", "ox"))   // 123
	fmt.Println(strings.TrimPrefix("oxo123", "ox")) // o123
	fmt.Println(strings.Trim("oxo123oxo", "ox"))    // 123
}
