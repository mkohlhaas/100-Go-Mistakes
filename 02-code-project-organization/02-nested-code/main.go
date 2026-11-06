package main

import "errors"

// join1 demonstrates the "arrow anti-pattern" (also called arrow code).
// Every validity check is wrapped inside an `else` block, so the successful
// path keeps drifting to the right — the code looks like the head of an arrow.
//
// This style is hard to read: the happy path is buried at the deepest
// indentation level, and adding another check means wrapping everything
// in yet another `else`.
func join1(s1, s2 string, max int) (string, error) {
	if s1 == "" {
		return "", errors.New("s1 is empty")
	} else {
		if s2 == "" {
			return "", errors.New("s2 is empty")
		} else {
			concat, err := concatenate(s1, s2)
			if err != nil {
				return "", err
			} else {
				if len(concat) > max {
					return concat[:max], nil
				} else {
					return concat, nil
				}
			}
		}
	}
}

// join2 demonstrates the idiomatic Go style: guard clauses with early returns.
// Each invalid condition returns immediately, so the remaining code is the
// linear happy path at the left margin.
//
// Benefits over join1:
//   - Flat structure — no deep indentation
//   - Happy path reads top-to-bottom without jumping through `else` blocks
//   - Easy to extend: add another check as one more guard clause
func join2(s1, s2 string, max int) (string, error) {
	if s1 == "" {
		return "", errors.New("s1 is empty")
	}
	if s2 == "" {
		return "", errors.New("s2 is empty")
	}
	concat, err := concatenate(s1, s2)
	if err != nil {
		return "", err
	}
	if len(concat) > max {
		return concat[:max], nil
	}
	return concat, nil
}

// concatenate is a stub used by both join functions.
// In real code it would combine s1 and s2 and return an error on failure.
func concatenate(_, _ string) (string, error) {
	return "", nil
}
