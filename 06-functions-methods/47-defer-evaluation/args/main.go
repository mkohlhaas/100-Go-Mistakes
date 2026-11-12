package main

import "fmt"

const (
	StatusSuccess  = "success"
	StatusErrorFoo = "error_foo"
	StatusErrorBar = "error_bar"
)

func main() {
	_ = f1()
	_ = f2()
	_ = f3()
}

// f1 demonstrates the mistake: defer captures status BY VALUE ("") at the
// time the defer statements run, so reassigning status below has no effect on
// what notify/incrementCounter print. These assignments are intentionally
// ineffectual — they illustrate the bug.
func f1() error {
	var status string
	defer notify(status)
	defer incrementCounter(status)

	if err := foo(); err != nil {
		status = StatusErrorFoo
		return err
	}

	if err := bar(); err != nil {
		status = StatusErrorBar
		return err
	}

	status = StatusSuccess
	return nil
}

func f2() error {
	var status string
	defer notifyPtr(&status)
	defer incrementCounterPtr(&status)

	if err := foo(); err != nil {
		status = StatusErrorFoo
		return err
	}

	if err := bar(); err != nil {
		status = StatusErrorBar
		return err
	}

	status = StatusSuccess
	return nil
}

func f3() error {
	var status string
	defer func() {
		notify(status)
		incrementCounter(status)
	}()

	if err := foo(); err != nil {
		status = StatusErrorFoo
		return err
	}

	if err := bar(); err != nil {
		status = StatusErrorBar
		return err
	}

	status = StatusSuccess
	return nil
}

func notify(status string) {
	fmt.Println("notify:", status)
}

func incrementCounter(status string) {
	fmt.Println("increment:", status)
}

func notifyPtr(status *string) {
	fmt.Println("notify:", *status)
}

func incrementCounterPtr(status *string) {
	fmt.Println("increment:", *status)
}

func foo() error {
	return nil
}

func bar() error {
	return nil
}
