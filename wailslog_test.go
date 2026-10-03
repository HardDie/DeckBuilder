package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
)

func TestWailsLogger(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo})))
	t.Cleanup(func() { slog.SetDefault(prev) })

	wailsLogger{}.Debug("hidden at Info")
	wailsLogger{}.Info("window ready")
	wailsLogger{}.Warning("slow")
	wailsLogger{}.Error("bridge failed") // the source must be this line, not wailslog.go

	var got []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("not JSON: %q", line)
		}
		got = append(got, m)
	}
	if len(got) != 3 {
		t.Fatalf("lines %v", got)
	}
	for i, want := range []string{"INFO", "WARN", "ERROR"} {
		if got[i]["level"] != want || got[i]["from"] != "wails" {
			t.Errorf("line %d: %v, want level %s from wails", i, got[i], want)
		}
	}
	src := got[2]["source"].(map[string]any)
	if filepath.Base(src["file"].(string)) != "wailslog_test.go" {
		t.Fatalf("source %v, want the caller", src)
	}
}
