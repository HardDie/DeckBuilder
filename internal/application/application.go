package application

import (
	"fmt"
	"net"
	"net/http"

	"github.com/HardDie/fsentry"
	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/api"
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
	repositoriesCard "github.com/HardDie/DeckBuilder/internal/repositories/card"
	repositoriesCollection "github.com/HardDie/DeckBuilder/internal/repositories/collection"
	repositoriesCore "github.com/HardDie/DeckBuilder/internal/repositories/core"
	repositoriesDeck "github.com/HardDie/DeckBuilder/internal/repositories/deck"
	repositoriesGame "github.com/HardDie/DeckBuilder/internal/repositories/game"
	repositoriesSettings "github.com/HardDie/DeckBuilder/internal/repositories/settings"
	serversCard "github.com/HardDie/DeckBuilder/internal/servers/card"
	serversCollection "github.com/HardDie/DeckBuilder/internal/servers/collection"
	serversDeck "github.com/HardDie/DeckBuilder/internal/servers/deck"
	serversGame "github.com/HardDie/DeckBuilder/internal/servers/game"
	serversGenerator "github.com/HardDie/DeckBuilder/internal/servers/generator"
	serversImage "github.com/HardDie/DeckBuilder/internal/servers/image"
	serversReplace "github.com/HardDie/DeckBuilder/internal/servers/replace"
	serversSearch "github.com/HardDie/DeckBuilder/internal/servers/search"
	serversSystem "github.com/HardDie/DeckBuilder/internal/servers/system"
	serversTTS "github.com/HardDie/DeckBuilder/internal/servers/tts"
	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
	servicesGenerator "github.com/HardDie/DeckBuilder/internal/services/generator"
	servicesReplace "github.com/HardDie/DeckBuilder/internal/services/replace"
	servicesSearch "github.com/HardDie/DeckBuilder/internal/services/search"
	servicesSystem "github.com/HardDie/DeckBuilder/internal/services/system"
	servicesTTS "github.com/HardDie/DeckBuilder/internal/services/tts"
)

type Application struct {
	cfg    *config.Config
	router *mux.Router
	tts    servicesTTS.TTS
}

func Get(debugFlag bool, version string) (*Application, error) {
	cfg := config.Get(debugFlag, version)

	routes := mux.NewRouter().StrictSlash(false)

	// static files
	api.RegisterStaticServer(routes)

	db := fsentry.New(cfg.Data, fsentry.WithPretty())
	if err := db.Init(); err != nil {
		logger.Error.Fatal(err)
	}

	core := repositoriesCore.New(db)

	err := core.Init()
	if err != nil {
		logger.Error.Fatal(err)
	}

	// system
	repositorySettings := repositoriesSettings.New(cfg, db)
	serviceSystem := servicesSystem.New(repositorySettings)
	serverSystem := serversSystem.New(cfg, serviceSystem)
	api.RegisterSystemServer(routes, serverSystem)

	// game
	repositoryGame := repositoriesGame.New(cfg, db)
	serviceGame := servicesGame.New(cfg, repositoryGame)
	serverGame := serversGame.New(*cfg, serviceGame, serverSystem)
	api.RegisterGameServer(routes, serverGame)

	// collection
	repositoryCollection := repositoriesCollection.New(cfg, db)
	serviceCollection := servicesCollection.New(cfg, repositoryCollection)
	serverCollection := serversCollection.New(*cfg, serviceCollection, serverSystem)
	api.RegisterCollectionServer(routes, serverCollection)

	// deck
	repositoryDeck := repositoriesDeck.New(cfg, db)
	serviceDeck := servicesDeck.New(cfg, repositoryDeck)
	serverDeck := serversDeck.New(*cfg, serviceDeck, serverSystem)
	api.RegisterDeckServer(routes, serverDeck)

	// card
	repositoryCard := repositoriesCard.New(cfg, db)
	serviceCard := servicesCard.New(cfg, repositoryCard)
	serverCard := serversCard.New(*cfg, serviceCard, serverSystem)
	api.RegisterCardServer(routes, serverCard)

	// image
	serverImage := serversImage.New(serviceGame, serviceCollection, serviceDeck, serviceCard)
	api.RegisterImageServer(routes, serverImage)

	// tts service
	serviceTTS := servicesTTS.New()
	serverTTS := serversTTS.New(serviceTTS)
	api.RegisterTTSServer(routes, serverTTS)

	// generator
	serviceGenerator := servicesGenerator.New(cfg, serviceGame, serviceCollection, serviceDeck, serviceCard, serviceSystem, serviceTTS)
	serverGenerator := serversGenerator.New(serviceGenerator)
	api.RegisterGeneratorServer(routes, serverGenerator)

	// replace
	serviceReplace := servicesReplace.New(serviceTTS)
	serverReplace := serversReplace.New(serviceReplace)
	api.RegisterReplaceServer(routes, serverReplace)

	// recursive search
	serviceSearch := servicesSearch.New(serviceGame, serviceCollection, serviceDeck, serviceCard)
	serverSearch := serversSearch.New(serviceSearch)
	api.RegisterSearchServer(routes, serverSearch)

	routes.Use(corsMiddleware)
	return &Application{
		cfg:    cfg,
		router: routes,
		tts:    serviceTTS,
	}, nil
}

func (app *Application) Handler() http.Handler {
	return app.router
}

func (app *Application) Listen() (net.Listener, int, error) {
	ln, port, err := listenLoopback(config.HTTPHost, config.HTTPPort, config.HTTPPortAttempts)
	if err != nil {
		return nil, 0, err
	}
	app.tts.SetHTTPPort(port)
	logger.Info.Printf("Listening on %s:%d...", config.HTTPHost, port)
	return ln, port, nil
}

func (app *Application) Serve(ln net.Listener) error {
	return http.Serve(ln, app.router)
}

func (app *Application) Run() error {
	ln, port, err := app.Listen()
	if err != nil {
		return err
	}
	if !app.cfg.Debug {
		network.OpenBrowser(fmt.Sprintf("http://%s:%d", config.HTTPHost, port))
	}
	return app.Serve(ln)
}

// CORS headers
func corsSetupHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, ContentType")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsSetupHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
