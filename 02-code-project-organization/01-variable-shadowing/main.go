package main

import (
	"net/http"
)

// tracing is a package-level flag simulating whether tracing is enabled.
var tracing bool

// createClientWithTracing is a placeholder that would return an instrumented client.
func createClientWithTracing() (*http.Client, error) {
	return nil, nil
}

// createDefaultClient is a placeholder that would return a plain client.
func createDefaultClient() (*http.Client, error) {
	return nil, nil
}

func main() {
	_ = listing1() // BUG: outer client stays nil (variable shadowing)
	_ = listing2() // fix: assign through throwaway `c`, then `client = c`
	_ = listing3() // fix: declare both vars outside, use `=` inside branches
	_ = listing4() // fix: cleanest — declare up front, assign with `=`, check err once
}

// listing1 demonstrates variable shadowing: the outer `client` is never assigned.
//
// `client` is declared in the function body (the outer block). Inside each branch
// of the `if/else`, `client, err := ...` uses the short variable declaration `:=`.
// Because `client` was NOT declared in the same block as those `:=` statements,
// Go creates a NEW `client` (and a NEW `err`) scoped to that branch only. These
// inner variables shadow the outer one and are discarded when the branch ends.
//
// As a result, the outer `client` keeps its zero value (nil) and the reference at
// the bottom of the function points to a nil client — even though a client was
// successfully created inside the branch.
//
// Compiles cleanly with `go vet`; this is what makes shadowing dangerous.
func listing1() error {
	var client *http.Client
	if tracing {
		client, err := createClientWithTracing() // shadows outer `client` + new `err`
		if err != nil {
			return err
		}
		_ = client
	} else {
		client, err := createDefaultClient() // shadows outer `client` + new `err`
		if err != nil {
			return err
		}
		_ = client
	}

	_ = client // outer client is still nil!
	return nil
}

// listing2 fixes shadowing by assigning to the outer `client` via a throwaway.
//
// `client` is declared in the outer scope. Inside each branch we declare a local
// `c` with `:=` (this is fine — `c` is new to the branch) and then copy it into
// the outer `client` with a simple assignment `client = c`. Simple assignment
// does not create a new variable; it writes to the existing outer `client`.
func listing2() error {
	var client *http.Client
	if tracing {
		c, err := createClientWithTracing() // local `c`, new `err`
		if err != nil {
			return err
		}
		client = c // assign to OUTER client
	} else {
		c, err := createDefaultClient() // local `c`, new `err`
		if err != nil {
			return err
		}
		client = c // assign to OUTER client
	}

	_ = client
	return nil
}

// listing3 fixes shadowing by declaring both variables in the outer scope.
//
// `client` and `err` are both declared up front in the function body. Inside the
// branches we use `=` (simple assignment), not `:=`, so both names refer to the
// outer variables. Each branch reassigns the same outer `err`, which is fine as
// long as we check it before it gets overwritten by the next branch (not an issue
// here since we return on error immediately).
func listing3() error {
	var client *http.Client
	var err error
	if tracing {
		client, err = createClientWithTracing() // assign to OUTER vars
		if err != nil {
			return err
		}
	} else {
		client, err = createDefaultClient() // assign to OUTER vars
		if err != nil {
			return err
		}
	}

	_ = client
	return nil
}

// listing4 is the cleanest fix: declare up front, assign with `=`, check err once.
//
// Declaring both `client` and `err` outside the `if/else` and using simple
// assignment inside keeps a single scope for both variables. Because neither
// branch returns early on error, we check `err` once after the `if/else`
// completes.
func listing4() error {
	var client *http.Client
	var err error
	if tracing {
		client, err = createClientWithTracing() // assign to OUTER vars
	} else {
		client, err = createDefaultClient() // assign to OUTER vars
	}
	if err != nil {
		return err
	}

	_ = client
	return nil
}
