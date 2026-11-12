package main

import (
	"fmt"
	"os"
)

// readFiles1 demonstrates the mistake: defer inside a loop.
//
// Each defer file.Close() is registered but does not run until readFiles1
// returns, so every file stays open for the whole loop. With many paths this
// leaks file descriptors and can hit the OS open-file limit. The deferred
// closures also all share the same `file` variable, so they end up closing
// only the last opened file repeatedly.
func readFiles1(ch <-chan string) (err error) {
	for path := range ch {
		var file *os.File
		file, err = os.Open(path)
		if err != nil {
			return err
		}

		defer func() {
			if cerr := file.Close(); cerr != nil {
				err = cerr
			}
		}()

		// Do something with file
	}
	return nil
}

// with helper function
func readFiles2(ch <-chan string) error {
	for path := range ch {
		if err := readFile(path); err != nil {
			return err
		}
	}
	return nil
}

// the helper function
func readFile(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()

	// Do something with file
	return nil
}

// with immediately-invoked closure
func readFiles3(ch <-chan string) error {
	for path := range ch {
		// the closure
		err := func() (err error) {
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer func() {
				if cerr := file.Close(); cerr != nil {
					err = cerr
				}
			}()

			// Do something with file
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	paths := []string{
		"/tmp/opencode/defer-loop-1.txt",
		"/tmp/opencode/defer-loop-2.txt",
		"/tmp/opencode/defer-loop-3.txt",
	}
	for _, p := range paths {
		if err := os.WriteFile(p, []byte("hello\n"), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", p, err)
			os.Exit(1)
		}
		defer func() {
			if err := os.Remove(p); err != nil {
				fmt.Fprintf(os.Stderr, "remove %s: %v\n", p, err)
			}
		}()
	}

	run := func(name string, fn func(<-chan string) error) {
		ch := make(chan string, len(paths))
		for _, p := range paths {
			ch <- p
		}
		close(ch)

		if err := fn(ch); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
			os.Exit(1)
		}
		fmt.Printf("%s completed\n", name)
	}

	run("readFiles1", readFiles1)
	run("readFiles2", readFiles2)
	run("readFiles3", readFiles3)
}
