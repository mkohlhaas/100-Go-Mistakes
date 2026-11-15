package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerBuggy(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/buggy?err=1", nil)
	rec := httptest.NewRecorder()

	handler1(rec, req)

	res := rec.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("failed to close response body: %v", err)
		}
	}()
	body, _ := io.ReadAll(res.Body)

	// The bug: the response body contains both the error and the success payload.
	if got := string(body); !strings.Contains(got, "foo") || !strings.Contains(got, "all good") {
		t.Fatalf("expected corrupted body mixing error and success, got %q", got)
	}
	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.StatusCode)
	}
}

func TestHandlerFixed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/fixed?err=1", nil)
	rec := httptest.NewRecorder()

	handler2(rec, req)

	res := rec.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("failed to close response body: %v", err)
		}
	}()
	body, _ := io.ReadAll(res.Body)

	// The fix: only the error body is present, and the status is 500.
	if got := string(body); strings.Contains(got, "all good") {
		t.Fatalf("expected no success payload after error, got %q", got)
	}
	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.StatusCode)
	}
}
