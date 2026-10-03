package main

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// wailsLogger sends the Wails framework's own messages to the default slog logger,
// so they reach the app log file (pkg/logger) like every other line.
// It satisfies Wails' logger.Logger. The recorded source is the Wails code that logged.
type wailsLogger struct{}

func (wailsLogger) Print(message string)   { emitWails(slog.LevelInfo, message) }
func (wailsLogger) Trace(message string)   { emitWails(slog.LevelDebug, message) }
func (wailsLogger) Debug(message string)   { emitWails(slog.LevelDebug, message) }
func (wailsLogger) Info(message string)    { emitWails(slog.LevelInfo, message) }
func (wailsLogger) Warning(message string) { emitWails(slog.LevelWarn, message) }
func (wailsLogger) Error(message string)   { emitWails(slog.LevelError, message) }

// Fatal logs the message and exits, as Wails expects.
func (wailsLogger) Fatal(message string) {
	emitWails(slog.LevelError, message)
	os.Exit(1)
}

func emitWails(level slog.Level, message string) {
	l := slog.Default()
	ctx := context.Background()
	if !l.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip Callers, emitWails, and the wailsLogger method
	r := slog.NewRecord(time.Now(), level, message, pcs[0])
	r.AddAttrs(slog.String("from", "wails"))
	_ = l.Handler().Handle(ctx, r)
}
