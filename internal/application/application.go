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
	"github.com/HardDie/DeckBuilder/internal/server"
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
	systemServer := server.NewSystemServer(cfg, serviceSystem)
	api.RegisterSystemServer(routes, systemServer)

	// game
	repositoryGame := repositoriesGame.New(cfg, game)
	serviceGame := servicesGame.New(cfg, repositoryGame)
	api.RegisterGameServer(routes, server.NewGameServer(serviceGame, systemServer))

	// collection
	repositoryCollection := repositoriesCollection.New(cfg, collection)
	serviceCollection := servicesCollection.New(cfg, repositoryCollection)
	api.RegisterCollectionServer(routes, server.NewCollectionServer(serviceCollection, systemServer))

	// deck
	repositoryDeck := repositoriesDeck.New(cfg, collection, deck)
	serviceDeck := servicesDeck.New(cfg, repositoryDeck)
	api.RegisterDeckServer(routes, server.NewDeckServer(serviceDeck, systemServer))

	// card
	repositoryCard := repositoriesCard.New(cfg, card)
	serviceCard := servicesCard.New(cfg, repositoryCard)
	api.RegisterCardServer(routes, server.NewCardServer(serviceCard, systemServer))

	// image
	api.RegisterImageServer(routes, server.NewImageServer(serviceGame, serviceCollection, serviceDeck, serviceCard))

	// tts service
	serviceTTS := servicesTTS.New()
	api.RegisterTTSServer(routes, server.NewTTSServer(serviceTTS))

	// generator
	serviceGenerator := servicesGenerator.New(cfg, serviceGame, serviceCollection, serviceDeck, serviceCard, serviceSystem, serviceTTS)
	api.RegisterGeneratorServer(routes, server.NewGeneratorServer(serviceGenerator))

	// replace
	serviceReplace := servicesReplace.New(serviceTTS)
	api.RegisterReplaceServer(routes, server.NewReplaceServer(serviceReplace))

	// recursive search
	serviceSearch := servicesSearch.New(serviceGame, serviceCollection, serviceDeck, serviceCard)
	api.RegisterSearchServer(routes, server.NewSearchServer(serviceSearch))

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
