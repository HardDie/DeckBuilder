package main

import (
	"context"
	"fmt"
	"net"
)

// App is the Wails bindings stub. Catalog and generate stay on the HTTP API.
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

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
