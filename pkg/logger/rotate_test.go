package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(path string) string {
	body, err := os.ReadFile(path)
	if err != nil {
		return "<missing>"
	}
	return string(body)
}

func TestRotate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	writeFile(t, path, "current, too big")
	writeFile(t, numbered(path, 1), "one")
	writeFile(t, numbered(path, 2), "two")
	writeFile(t, numbered(path, 3), "three")

	if err := rotate(path, 5, 3); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		path:              "<missing>",
		numbered(path, 1): "current, too big",
		numbered(path, 2): "one",
		numbered(path, 3): "two", // "three" was the oldest and is gone
	}
	for p, body := range want {
		if got := readFile(p); got != body {
			t.Errorf("%s: %q, want %q", filepath.Base(p), got, body)
		}
	}
	if filepath.Base(numbered(path, 1)) != "app.1.log" {
		t.Errorf("numbered name %s", filepath.Base(numbered(path, 1)))
	}
}

func TestRotateLeavesFile(t *testing.T) {
	dir := t.TempDir()
	small := filepath.Join(dir, "small.log")
	writeFile(t, small, "small")
	big := filepath.Join(dir, "big.log")
	writeFile(t, big, "big enough")

	tests := []struct {
		name          string
		path          string
		maxSize, keep int
	}{
		{"under_the_limit", small, 5, 3},
		{"no_file_yet", filepath.Join(dir, "none.log"), 5, 3},
		{"rotation_off_by_size", big, 0, 3},
		{"rotation_off_by_keep", big, 5, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := readFile(tt.path)
			if err := rotate(tt.path, int64(tt.maxSize), tt.keep); err != nil {
				t.Fatal(err)
			}
			if got := readFile(tt.path); got != before {
				t.Fatalf("file changed: %q → %q", before, got)
			}
		})
	}
}
