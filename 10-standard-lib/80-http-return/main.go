package main

import (
	"fmt"
	"net/http"
)

func handler1(w http.ResponseWriter, req *http.Request) {
	err := foo(req)
	if err != nil {
		http.Error(w, "foo", http.StatusInternalServerError) // ⚠️ no immediate return
	}

	_, _ = w.Write([]byte("all good"))
	w.WriteHeader(http.StatusCreated)
}

func handler2(w http.ResponseWriter, req *http.Request) {
	err := foo(req)
	if err != nil {
		http.Error(w, "foo", http.StatusInternalServerError)
		return // ✅ immediate return
	}

	_, _ = w.Write([]byte("all good"))
	w.WriteHeader(http.StatusCreated)
}

func foo(req *http.Request) error {
	if req.URL.Query().Get("err") == "1" {
		return fmt.Errorf("simulated failure")
	}
	return nil
}
