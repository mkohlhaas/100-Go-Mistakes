package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func getBytes1(reader io.Reader) ([]byte, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return []byte(sanitize1(string(b))), nil
}

func sanitize1(s string) string {
	return strings.TrimSpace(s)
}

func getBytes2(reader io.Reader) ([]byte, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return sanitize2(b), nil
}

func sanitize2(b []byte) []byte {
	return bytes.TrimSpace(b)
}

func main() {
	input := "   hello, world!   "
	fmt.Printf("input: %q\n", input) // input: "   hello, world!   "

	out1, err := getBytes1(strings.NewReader(input))
	if err != nil {
		panic(err)
	}
	fmt.Printf("getBytes1: %q\n", out1) // getBytes1: "hello, world!"

	out2, err := getBytes2(strings.NewReader(input))
	if err != nil {
		panic(err)
	}
	fmt.Printf("getBytes2: %q\n", out2) // getBytes2: "hello, world!"
}
