package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type transientError struct {
	err error
}

func (t transientError) Error() string {
	return fmt.Sprintf("transient error: %v", t.err)
}

// handler1 uses a type switch — breaks when the error is wrapped.
func handler1(w http.ResponseWriter, r *http.Request) {
	transactionID := r.URL.Query().Get("transaction")

	amount, err := getTransactionAmount1(transactionID)
	if err != nil {
		switch err := err.(type) {
		case transientError:
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	_ = amount
}

func getTransactionAmount1(transactionID string) (float32, error) {
	if len(transactionID) != 5 {
		return 0, fmt.Errorf("id is invalid: %s", transactionID)
	}

	amount, err := getTransactionAmountFromDB1(transactionID)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction %s: %w", transactionID, err)
	}
	return amount, nil
}

func getTransactionAmountFromDB1(id string) (float32, error) {
	if id == "trans" {
		return 0, transientError{err: fmt.Errorf("db unavailable")}
	}
	return 0, nil
}

// handler2 uses errors.As — walks the wrapped chain correctly.
func handler2(w http.ResponseWriter, r *http.Request) {
	transactionID := r.URL.Query().Get("transaction")

	amount, err := getTransactionAmount2(transactionID)
	if err != nil {
		if errors.As(err, &transientError{}) {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	_ = amount
}

func getTransactionAmount2(transactionID string) (float32, error) {
	if len(transactionID) != 5 {
		return 0, fmt.Errorf("id is invalid: %s", transactionID)
	}

	amount, err := getTransactionAmountFromDB2(transactionID)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction %s: %w", transactionID, err)
	}
	return amount, nil
}

func getTransactionAmountFromDB2(id string) (float32, error) {
	if id == "trans" {
		return 0, transientError{err: fmt.Errorf("db unavailable")}
	}
	return 0, nil
}

func main() {
	cases := []struct {
		name  string
		id    string
		want1 int
		want2 int
	}{
		{"valid id", "abcde", http.StatusOK, http.StatusOK},
		{"invalid id", "abc", http.StatusBadRequest, http.StatusBadRequest},
		{"valid id but transient db error", "trans", http.StatusServiceUnavailable, http.StatusServiceUnavailable},
	}

	for _, c := range cases {
		fmt.Printf("=== %s (id=%q) ===\n", c.name, c.id)
		got1 := callHandler(handler1, c.id)
		got2 := callHandler(handler2, c.id)
		fmt.Printf("  handler1 (type switch): %d%s\n", got1, mark(got1, c.want1))
		fmt.Printf("  handler2 (errors.As):   %d%s\n", got2, mark(got2, c.want2))
	}
}

func callHandler(h http.HandlerFunc, id string) int {
	req := httptest.NewRequest("GET", "/?transaction="+id, nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec.Code
}

func mark(got, want int) string {
	if got == want {
		return " ✓"
	}
	return " ✗"
}
