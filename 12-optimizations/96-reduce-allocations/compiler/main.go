package main

import "fmt"

type cache struct {
	m map[string]int
}

//go:noinline
func (c *cache) get1(bytes []byte) (v int, contains bool) {
	key := string(bytes)
	v, contains = c.m[key]
	return
}

//go:noinline
func (c *cache) get2(bytes []byte) (v int, contains bool) {
	v, contains = c.m[string(bytes)]
	return
}

func main() {
	c := &cache{m: map[string]int{"hello": 1, "world": 2}}

	for _, key := range [][]byte{{'h', 'e', 'l', 'l', 'o'}, {'w', 'o', 'r', 'l', 'd'}} {
		v1, ok1 := c.get1(key)
		v2, ok2 := c.get2(key)
		fmt.Printf("key=%q get1=(%d, %v) get2=(%d, %v)\n", string(key), v1, ok1, v2, ok2)
	}
}
