package main

import (
	"errors"
	"fmt"
)

func bar() error {
	return barError{}
}

type barError struct{}

func (b barError) Error() string {
	return "bar error"
}

func listing1() error {
	err := bar()
	if err != nil {
		return err
	}
	// ...
	return nil
}

type BarError struct {
	Err error
}

func (b BarError) Error() string {
	return "bar failed:" + b.Err.Error()
}

func listing2() error {
	err := bar()
	if err != nil {
		return BarError{Err: err}
	}
	// ...
	return nil
}

func listing3() error {
	err := bar()
	if err != nil {
		return fmt.Errorf("bar failed: %w", err)
	}
	// ...
	return nil
}

func listing4() error {
	err := bar()
	if err != nil {
		return fmt.Errorf("bar failed: %v", err)
	}
	// ...
	return nil
}

// listing5 builds a longer chain: listing5 -> listing3 -> bar.
func listing5() error {
	err := listing3()
	if err != nil {
		return fmt.Errorf("outer operation failed: %w", err)
	}
	return nil
}

func main() {
	listings := []struct {
		name string
		fn   func() error
	}{
		{"listing1 (unwrapped)", listing1},
		{"listing2 (custom wrapper, no Unwrap)", listing2},
		{"listing3 (%w wraps error)", listing3},
		{"listing4 (%v loses error)", listing4},
		{"listing5 (longer chain)", listing5},
	}

	for _, l := range listings {
		err := l.fn()
		fmt.Printf("%s -> %v\n", l.name, err)
		if err == nil {
			continue
		}

		unwrapped := errors.Unwrap(err)
		fmt.Printf("  errors.Unwrap: %v\n", unwrapped)
		fmt.Printf("  errors.Is(barError{}): %v\n", errors.Is(err, barError{}))
	}
	fmt.Println()

	// Demonstrate that only listing3's wrapped error can be matched with errors.Is.
	err3 := listing3()
	if errors.Is(err3, barError{}) {
		fmt.Println("listing3: errors.Is matches barError{} (correct wrapping with %w)")
	}

	err4 := listing4()
	if !errors.Is(err4, barError{}) {
		fmt.Println("listing4: errors.Is does NOT match barError{} (lost the original error)")
	}

	// Walk a longer chain: listing5 -> listing3 -> bar.
	err5 := listing5()
	fmt.Println()
	fmt.Println("listing5 chain walk:")
	for err := err5; err != nil; err = errors.Unwrap(err) {
		fmt.Printf("  %v\n", err)
	}
	fmt.Printf("  errors.Is(err5, barError{}): %v\n", errors.Is(err5, barError{}))
}
