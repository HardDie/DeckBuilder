//go:build !nomain

// Wails desktop shell. The window hosts the Vue UI.
// Loopback HTTP serves images and TTS only.
package main

import (
	"context"
	"embed"
	"errors"
	"net/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/HardDie/DeckBuilder/bindings/card"
	"github.com/HardDie/DeckBuilder/bindings/collection"
	"github.com/HardDie/DeckBuilder/bindings/deck"
	"github.com/HardDie/DeckBuilder/bindings/game"
	bindingsGenerator "github.com/HardDie/DeckBuilder/bindings/generator"
	bindingsReplace "github.com/HardDie/DeckBuilder/bindings/replace"
	bindingsSearch "github.com/HardDie/DeckBuilder/bindings/search"
	bindingsSystem "github.com/HardDie/DeckBuilder/bindings/system"
	"github.com/HardDie/DeckBuilder/internal/application"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/servers"
	"github.com/HardDie/DeckBuilder/pkg/version"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	backend, err := application.Get(version.String())
	if err != nil {
		logger.Error.Fatal(err.Error())
	}

	httpServer := servers.New(
		backend.GameService(),
		backend.CollectionService(),
		backend.DeckService(),
		backend.CardService(),
		backend.TTSService(),
	)
	ln, _, err := httpServer.Listen()
	if err != nil {
		logger.Error.Fatal(err.Error())
	}
	go func() {
		if err := httpServer.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error.Println(err.Error())
		}
	}()

	cfg := *backend.Config()
	app := NewApp(ln)
	games := game.New(cfg, backend.GameService())
	collections := collection.New(cfg, backend.CollectionService())
	decks := deck.New(cfg, backend.DeckService())
	cards := card.New(cfg, backend.CardService())
	systems := bindingsSystem.New(cfg, backend.SystemService())
	searches := bindingsSearch.New(backend.SearchService())
	generators := bindingsGenerator.New(backend.GeneratorService())
	replaces := bindingsReplace.New(backend.ReplaceService())

	err = wails.Run(&options.App{
		Title:  "DeckBuilder",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: apiHandler(httpServer.Handler()),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			game.BindWindow(games, ctx)
		},
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			games,
			collections,
			decks,
			cards,
			systems,
			searches,
			generators,
			replaces,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
