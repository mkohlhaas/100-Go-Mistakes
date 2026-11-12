package main

import (
	"errors"
	"fmt"
	"strings"
)

func main() {
	// A substring shares the underlying byte array of the original string.
	// Keeping the substring keeps the whole string alive.
	s1 := "Hello, World!"
	s2 := s1[:5]
	fmt.Println(s2) // Hello

	// Slicing by bytes breaks multi-byte characters; convert to runes first.
	s1 = "Hêllo, World!"
	s2 = string([]rune(s1)[:5])
	fmt.Println(s2) // Hêllo

	// Real-world scenario: extract a fixed-length id from a log line and keep it.
	s := &store{}
	logs := []string{
		"00000000-0000-0000-0000-000000000001 some extra payload",
		"00000000-0000-0000-0000-000000000002 some extra payload",
		"short", // not correctly formatted
		"00000000-0000-0000-0000-000000000003 some extra payload",
	}
	for _, log := range logs {
		if err := s.handleLog1(log); err != nil {
			fmt.Println("handleLog1:", err)
		}
		if err := s.handleLog2(log); err != nil {
			fmt.Println("handleLog2:", err)
		}
		if err := s.handleLog3(log); err != nil {
			fmt.Println("handleLog3:", err)
		}
	}
	fmt.Println("stored UUIDs:", s.uuids)
}

type store struct {
	uuids []string
}

func (s *store) handleLog1(log string) error {
	if len(log) < 36 {
		return errors.New("log is not correctly formatted")
	}
	// BUG: uuid shares log's backing array, so the entire log string
	// stays in memory for as long as uuid is retained.
	uuid := log[:36]
	s.uuids = append(s.uuids, uuid)
	// Do something
	return nil
}

func (s *store) handleLog2(log string) error {
	if len(log) < 36 {
		return errors.New("log is not correctly formatted")
	}
	// Fixed: convert the slice back to a string via []byte, which forces a copy.
	uuid := string([]byte(log[:36]))
	s.uuids = append(s.uuids, uuid)
	// Do something
	return nil
}

func (s *store) handleLog3(log string) error {
	if len(log) < 36 {
		return errors.New("log is not correctly formatted")
	}
	// Fixed: strings.Clone (Go 1.17+) copies the substring explicitly.
	uuid := strings.Clone(log[:36])
	s.uuids = append(s.uuids, uuid)
	// Do something
	return nil
}
