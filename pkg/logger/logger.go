// Package logger sets up log/slog for an app: a JSON or text log on the console
// and, optionally, in a file that is rotated by size when it is opened.
//
// After Init, code logs through the standard slog functions:
//
//	slog.Info("render finished", "game", id, "sheets", n)
//	slog.Error("render failed", "err", err)
//
// Lines from the older log package go to the same place.
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

// Format is how each log line is written.
type Format string

const (
	// FormatJSON writes one JSON object per line. It is the default.
	FormatJSON Format = "json"
	// FormatText writes key=value pairs, easier to read in an editor.
	FormatText Format = "text"
)

// Options configures Init. The zero value logs JSON at Info to stderr only.
type Options struct {
	// Dir is the folder of the log file. Empty means console only.
	Dir string
	// File is the log file name inside Dir. Default "app.log".
	File string
	// Format of each line. Default FormatJSON.
	Format Format
	// Level is the lowest level written. Default slog.LevelInfo.
	Level slog.Level
	// MaxSize rotates the file when it is opened larger than this many bytes. 0 never rotates.
	MaxSize int64
	// Keep is how many rotated files stay: app.1.log … app.<Keep>.log.
	Keep int
	// Console also receives every line. Default os.Stderr; nil after Init means none.
	Console io.Writer
}

var (
	mu   sync.Mutex
	file *os.File
)

// Init makes the default slog logger write as opts says and returns the log file path.
// If the file cannot be used, the logger still writes to the console
// and Init returns the error, so the app can start anyway.
func Init(opts Options) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	closeFile()

	console := opts.Console
	if console == nil {
		console = os.Stderr
	}
	if opts.Dir == "" {
		slog.SetDefault(slog.New(newHandler(console, opts)))
		return "", nil
	}

	name := opts.File
	if name == "" {
		name = "app.log"
	}
	path := filepath.Join(opts.Dir, name)
	f, err := openRotated(path, opts.MaxSize, opts.Keep)
	if err != nil {
		slog.SetDefault(slog.New(newHandler(console, opts)))
		return "", err
	}
	file = f
	slog.SetDefault(slog.New(newHandler(io.MultiWriter(console, f), opts)))
	return path, nil
}

// Close closes the log file and leaves a console-only JSON logger.
func Close() error {
	mu.Lock()
	defer mu.Unlock()
	slog.SetDefault(slog.New(newHandler(os.Stderr, Options{})))
	return closeFile()
}

func closeFile() error {
	if file == nil {
		return nil
	}
	err := file.Close()
	file = nil
	return err
}

func openRotated(path string, maxSize int64, keep int) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create log folder: %w", err)
	}
	if err := rotate(path, maxSize, keep); err != nil {
		return nil, fmt.Errorf("rotate log: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open log: %w", err)
	}
	return f, nil
}

// newHandler adds the source file and line only to Warn and Error lines.
func newHandler(w io.Writer, opts Options) slog.Handler {
	build := func(source bool) slog.Handler {
		ho := &slog.HandlerOptions{Level: opts.Level, AddSource: source, ReplaceAttr: shortSource}
		if opts.Format == FormatText {
			return slog.NewTextHandler(w, ho)
		}
		return slog.NewJSONHandler(w, ho)
	}
	return sourceOnWarn{plain: build(false), withSource: build(true)}
}

// shortSource keeps only the file name of the source; the function already names the package.
func shortSource(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		if src, ok := a.Value.Any().(*slog.Source); ok {
			src.File = filepath.Base(src.File)
		}
	}
	return a
}

// sourceOnWarn sends Warn and above to a handler that records the source.
type sourceOnWarn struct {
	plain, withSource slog.Handler
}

func (h sourceOnWarn) Enabled(ctx context.Context, l slog.Level) bool {
	return h.plain.Enabled(ctx, l)
}

func (h sourceOnWarn) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= slog.LevelWarn {
		return h.withSource.Handle(ctx, r)
	}
	return h.plain.Handle(ctx, r)
}

func (h sourceOnWarn) WithAttrs(attrs []slog.Attr) slog.Handler {
	return sourceOnWarn{plain: h.plain.WithAttrs(attrs), withSource: h.withSource.WithAttrs(attrs)}
}

func (h sourceOnWarn) WithGroup(name string) slog.Handler {
	return sourceOnWarn{plain: h.plain.WithGroup(name), withSource: h.withSource.WithGroup(name)}
}
