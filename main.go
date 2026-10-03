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
	"github.com/HardDie/DeckBuilder/bindings/errfmt"
	"github.com/HardDie/DeckBuilder/bindings/game"
	bindingsGenerator "github.com/HardDie/DeckBuilder/bindings/generator"
	bindingsReplace "github.com/HardDie/DeckBuilder/bindings/replace"
	bindingsSearch "github.com/HardDie/DeckBuilder/bindings/search"
	bindingsSystem "github.com/HardDie/DeckBuilder/bindings/system"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/servers"
	"github.com/HardDie/DeckBuilder/pkg/version"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	backend := wire(version.String())

	httpServer := servers.New(
		backend.game,
		backend.collection,
		backend.deck,
		backend.card,
		backend.tts,
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

	cfg := *backend.cfg
	app := NewApp(ln)
	games := game.New(cfg, backend.game)
	collections := collection.New(cfg, backend.collection)
	decks := deck.New(cfg, backend.deck)
	cards := card.New(cfg, backend.card)
	systems := bindingsSystem.New(cfg, backend.system)
	searches := bindingsSearch.New(backend.search)
	generators := bindingsGenerator.New(backend.generator)
	replaces := bindingsReplace.New(backend.replace)

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
		// Errors from bound methods reach the window as readable messages.
		ErrorFormatter: errfmt.Format,
		// One instance per user. A second launch exits and focuses the first window.
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "com.harddie.deckbuilder",
			OnSecondInstanceLaunch: app.focus,
		},
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
