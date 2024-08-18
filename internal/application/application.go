package application

import (
	"net/http"

	"github.com/HardDie/fsentry"
	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/api"
	"github.com/HardDie/DeckBuilder/internal/config"
	dbCard "github.com/HardDie/DeckBuilder/internal/db/card"
	dbCollection "github.com/HardDie/DeckBuilder/internal/db/collection"
	dbCore "github.com/HardDie/DeckBuilder/internal/db/core"
	dbDeck "github.com/HardDie/DeckBuilder/internal/db/deck"
	dbGame "github.com/HardDie/DeckBuilder/internal/db/game"
	dbSettings "github.com/HardDie/DeckBuilder/internal/db/settings"
	"github.com/HardDie/DeckBuilder/internal/logger"
	repositoriesCard "github.com/HardDie/DeckBuilder/internal/repositories/card"
	repositoriesCollection "github.com/HardDie/DeckBuilder/internal/repositories/collection"
	repositoriesDeck "github.com/HardDie/DeckBuilder/internal/repositories/deck"
	repositoriesGame "github.com/HardDie/DeckBuilder/internal/repositories/game"
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
	router *mux.Router
}

func Get(debugFlag bool, version string) (*Application, error) {
	cfg := config.Get(debugFlag, version)

	routes := mux.NewRouter().StrictSlash(false)

	// static files
	api.RegisterStaticServer(routes)

	// fsentry db
	fs := fsentry.NewFSEntry(cfg.Data, fsentry.WithPretty())

	// db methods
	core := dbCore.New(fs)
	settings := dbSettings.New(fs)
	game := dbGame.New(fs)
	collection := dbCollection.New(fs, game)
	deck := dbDeck.New(fs, collection)
	card := dbCard.New(fs, deck)

	err := core.Init()
	if err != nil {
		logger.Error.Fatal(err)
	}

	// system
	serviceSystem := servicesSystem.New(cfg, settings)
	serverSystem := serversSystem.New(cfg, serviceSystem)
	api.RegisterSystemServer(routes, serverSystem)

	// game
	repositoryGame := repositoriesGame.New(cfg, game)
	serviceGame := servicesGame.New(cfg, repositoryGame)
	serverGame := serversGame.New(*cfg, serviceGame, serverSystem)
	api.RegisterGameServer(routes, serverGame)

	// collection
	repositoryCollection := repositoriesCollection.New(cfg, collection)
	serviceCollection := servicesCollection.New(cfg, repositoryCollection)
	serverCollection := serversCollection.New(serviceCollection, serverSystem)
	api.RegisterCollectionServer(routes, serverCollection)

	// deck
	repositoryDeck := repositoriesDeck.New(cfg, collection, deck)
	serviceDeck := servicesDeck.New(cfg, repositoryDeck)
	serverDeck := serversDeck.New(serviceDeck, serverSystem)
	api.RegisterDeckServer(routes, serverDeck)

	// card
	repositoryCard := repositoriesCard.New(cfg, card)
	serviceCard := servicesCard.New(cfg, repositoryCard)
	serverCard := serversCard.New(serviceCard, serverSystem)
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
		router: routes,
	}, nil
}

func (app *Application) Run() error {
	http.Handle("/", app.router)
	logger.Info.Println("Listening on :5000...")
	return http.ListenAndServe("127.0.0.1:5000", nil)
}

// CORS headers
func corsSetupHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// CORS Headers middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsSetupHeaders(w)
		next.ServeHTTP(w, r)
	})
}
