package main

import (
	"context"
	"fmt"
)

// App is the Wails bindings stub. Catalog and generate stay on the HTTP API;
// do not add REST replacements here until the shell loads 127.0.0.1:<port>.
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
