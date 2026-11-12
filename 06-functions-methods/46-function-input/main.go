package main

import (
	"bufio"
	"io"
	"os"
)

func countEmptyLinesInFile(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}

	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		if scanner.Text() == "" {
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func countEmptyLines(reader io.Reader) (int, error) {
	scanner := bufio.NewScanner(reader)
	count := 0
	for scanner.Scan() {
		if scanner.Text() == "" {
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

func main() {
	file, err := os.Open("main.go")
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	_, _ = countEmptyLines(file)
}
