// Wails desktop shell: same HTTP application as cmd/deck_builder, plus a webview.
// Catalog stays on /api (see docs/architecture/010-wails-shell.md).
package main

import (
	"embed"
	"errors"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/HardDie/DeckBuilder/internal/application"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	backend, err := application.Get(os.Getenv("frontenddevserverurl") != "", version())
	if err != nil {
		logger.Error.Fatal(err.Error())
	}

	ln, _, err := backend.Listen()
	if err != nil {
		logger.Error.Fatal(err.Error())
	}
	go func() {
		if err := backend.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error.Println(err.Error())
		}
	}()

	app := NewApp(ln)

	err = wails.Run(&options.App{
		Title:  "DeckBuilder",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: apiHandler(backend.Handler()),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func version() string {
	if info, available := debug.ReadBuildInfo(); available {
		switch info.Main.Version {
		case "", "(devel)":
		default:
			return info.Main.Version
		}
	}
	return "unknown"
}
