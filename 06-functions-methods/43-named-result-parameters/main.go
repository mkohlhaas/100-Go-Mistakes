package main

import (
	"fmt"
	"io"
	"strings"
)

// f returns its argument via a named result parameter.
// The bare `return` statement returns the named result b.
func f(a int) (b int) {
	b = a
	return
}

// locator resolves addresses to coordinates.
type locator interface {
	getCoordinates(address string) (lat, lng float32, err error)
}

type loc struct{}

func (l loc) getCoordinates(address string) (lat, lng float32, err error) {
	return 0, 0, nil
}

// ReadFull reads from r into buf until buf is full or an error occurs.
// Named results n and err are accumulated and returned by the bare return.
func ReadFull(r io.Reader, buf []byte) (n int, err error) {
	for len(buf) > 0 && err == nil {
		var nr int
		nr, err = r.Read(buf)
		n += nr
		buf = buf[nr:]
	}
	return
}

// addNames documents its intent by naming the result sum.
func addNames(a, b int) (sum int) {
	sum = a + b
	return
}

// doublePitfall demonstrates the named-result shadowing pitfall.
// Inside the inner block, `result := x * 2` declares a NEW local variable
// named result that shadows the named result parameter. The bare return
// is outside that block, so it returns the named result (still 0),
// not the local variable. doublePitfall(5) returns 0 instead of 10.
func doublePitfall(x int) (result int) {
	{
		result := x * 2 // new local shadows named result
		_ = result
	}
	return // named result result is still 0
}

// doubleFixed correctly assigns to the named result without redeclaring it.
func doubleFixed(x int) (result int) {
	result = x * 2 // assigns to named result, no new variable
	return
}

func main() {
	fmt.Println("f(42) =\t", f(42))

	var l locator = loc{}
	lat, lng, err := l.getCoordinates("1600 Amphitheatre Parkway")
	fmt.Printf("getCoordinates: lat=%v lng=%v err=%v\n", lat, lng, err)

	reader := strings.NewReader("hello world")
	buf := make([]byte, 5)
	n, err := ReadFull(reader, buf)
	fmt.Printf("ReadFull: n=%d err=%v buf=%q\n", n, err, buf[:n])

	fmt.Printf("addNames(3, 4) =\t%d\n", addNames(3, 4))

	res := doublePitfall(5)
	fmt.Printf("doublePitfall(5) =\t%d (want 10)\n", res) // doublePitfall(5) =	0 (want 10)

	res = doubleFixed(5)
	fmt.Printf("doubleFixed(5) =\t%d\n", res) // doubleFixed(5) =	10
}
