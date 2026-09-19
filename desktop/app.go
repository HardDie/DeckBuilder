package main

import (
	"context"
	"net"
)

// App owns the Wails window lifecycle. Catalog verbs live in desktop/bindings/*.
type App struct {
	ctx context.Context
	ln  net.Listener
}

func NewApp(ln net.Listener) *App {
	return &App{ln: ln}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(_ context.Context) {
	if a.ln != nil {
		_ = a.ln.Close()
	}
}
