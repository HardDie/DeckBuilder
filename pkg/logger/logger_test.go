package logger

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// lines parses the JSON log file into one map per line.
func lines(t *testing.T, path string) []map[string]any {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var out []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(string(body)), "\n") {
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("not JSON: %q", line)
		}
		out = append(out, m)
	}
	return out
}

func initFile(t *testing.T, opts Options) (string, *bytes.Buffer) {
	t.Helper()
	console := &bytes.Buffer{}
	opts.Dir = filepath.Join(t.TempDir(), "logs")
	opts.Console = console
	path, err := Init(opts)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = Close() })
	return path, console
}

func TestJSONFileAndConsole(t *testing.T) {
	path, console := initFile(t, Options{File: "test.log"})
	if filepath.Base(path) != "test.log" {
		t.Fatalf("path %s", path)
	}

	slog.Info("render finished", "game", "raid", "sheets", 3)
	slog.Warn("slow download", "url", "x")
	slog.Debug("hidden at Info")

	got := lines(t, path)
	if len(got) != 2 {
		t.Fatalf("lines %v", got)
	}
	info, warn := got[0], got[1]
	if info["msg"] != "render finished" || info["level"] != "INFO" || info["game"] != "raid" || info["sheets"] != float64(3) {
		t.Fatalf("info line %v", info)
	}
	if _, ok := info["source"]; ok {
		t.Fatal("info lines carry no source")
	}
	src, ok := warn["source"].(map[string]any)
	if warn["level"] != "WARN" || !ok || src["file"] != "logger_test.go" {
		t.Fatalf("warn line needs its source: %v", warn)
	}
	if !strings.Contains(console.String(), `"msg":"render finished"`) {
		t.Fatalf("console did not get the line: %s", console.String())
	}
}

func TestTextFormatAndLevel(t *testing.T) {
	path, _ := initFile(t, Options{Format: FormatText, Level: slog.LevelDebug})
	if filepath.Base(path) != "app.log" {
		t.Fatalf("default file name %s", path)
	}
	slog.Debug("shown at Debug", "k", 1)

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `level=DEBUG msg="shown at Debug" k=1`) {
		t.Fatalf("text line %q", body)
	}
}

func TestStdLogGoesToFile(t *testing.T) {
	path, _ := initFile(t, Options{})
	log.Print("from the log package")
	if got := lines(t, path); len(got) != 1 || got[0]["msg"] != "from the log package" {
		t.Fatalf("lines %v", got)
	}
}

func TestCloseStopsTheFile(t *testing.T) {
	path, _ := initFile(t, Options{})
	slog.Info("before")
	if err := Close(); err != nil {
		t.Fatal(err)
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(&bytes.Buffer{}, nil))) // keep test output quiet
	slog.Info("after")
	if got := lines(t, path); len(got) != 1 {
		t.Fatalf("lines after Close %v", got)
	}
}

func TestConsoleOnly(t *testing.T) {
	console := &bytes.Buffer{}
	path, err := Init(Options{Console: console})
	t.Cleanup(func() { _ = Close() })
	if err != nil || path != "" {
		t.Fatalf("path %q err %v", path, err)
	}
	slog.Info("console only")
	if !strings.Contains(console.String(), "console only") {
		t.Fatal("console missed the line")
	}
}

func TestUnwritableFolderFallsBackToConsole(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("folder permissions differ on Windows")
	}
	parent := t.TempDir()
	if err := os.Chmod(parent, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(parent, 0o700) })

	console := &bytes.Buffer{}
	_, err := Init(Options{Dir: filepath.Join(parent, "logs"), Console: console})
	t.Cleanup(func() { _ = Close() })
	if err == nil {
		t.Fatal("want an error for an unwritable folder")
	}
	slog.Info("still logged")
	if !strings.Contains(console.String(), "still logged") {
		t.Fatal("the console logger should still work")
	}
}

func TestSetLevel(t *testing.T) {
	path, _ := initFile(t, Options{})
	slog.Debug("hidden before")
	SetLevel(slog.LevelDebug)
	slog.Debug("shown after")
	SetLevel(slog.LevelInfo)
	slog.Debug("hidden again")

	got := lines(t, path)
	if len(got) != 1 || got[0]["msg"] != "shown after" {
		t.Fatalf("lines %v", got)
	}
}
