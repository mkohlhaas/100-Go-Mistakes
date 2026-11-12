package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountEmptyLines(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:  "no empty lines",
			input: "foo\nbar\nbaz\n",
			want:  0,
		},
		{
			name:  "single empty line",
			input: "foo\n\nbar\n",
			want:  1,
		},
		{
			name:  "multiple empty lines",
			input: "foo\n\nbar\n\n\nbaz\n",
			want:  3,
		},
		{
			name:  "leading and trailing empty lines",
			input: "\nfoo\nbar\n\n",
			want:  2,
		},
		{
			name:  "empty input",
			input: "",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := countEmptyLines(strings.NewReader(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %d empty lines, want %d", got, tt.want)
			}
		})
	}
}

func TestCountEmptyLinesInFile(t *testing.T) {
	content := "foo\n\nbar\n\n\nbaz\n"

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	got, err := countEmptyLinesInFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3 {
		t.Errorf("got %d empty lines, want 3", got)
	}
}
