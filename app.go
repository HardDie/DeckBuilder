package main

import (
	"context"
	"net"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App owns the Wails window lifecycle. Catalog verbs live in bindings/*.
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

// focus brings the running window forward when a second launch is blocked.
func (a *App) focus(_ options.SecondInstanceData) {
	if a.ctx == nil {
		return
	}
	runtime.WindowUnminimise(a.ctx)
	runtime.WindowShow(a.ctx)
}
