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
	"github.com/HardDie/DeckBuilder/internal/service"
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
	systemService := service.NewService(cfg, settings)
	systemServer := server.NewSystemServer(cfg, systemService)
	api.RegisterSystemServer(routes, systemServer)

	// game
	repositoryGame := repositoriesGame.New(cfg, game)
	gameService := service.NewGameService(cfg, repositoryGame)
	api.RegisterGameServer(routes, server.NewGameServer(gameService, systemServer))

	// collection
	repositoryCollection := repositoriesCollection.New(cfg, collection)
	collectionService := service.NewCollectionService(cfg, repositoryCollection)
	api.RegisterCollectionServer(routes, server.NewCollectionServer(collectionService, systemServer))

	// deck
	repositoryDeck := repositoriesDeck.New(cfg, collection, deck)
	deckService := service.NewDeckService(cfg, repositoryDeck)
	api.RegisterDeckServer(routes, server.NewDeckServer(deckService, systemServer))

	// card
	repositoryCard := repositoriesCard.New(cfg, card)
	cardService := service.NewCardService(cfg, repositoryCard)
	api.RegisterCardServer(routes, server.NewCardServer(cardService, systemServer))

	// image
	api.RegisterImageServer(routes, server.NewImageServer(gameService, collectionService, deckService, cardService))

	// tts service
	ttsService := service.NewTTSService()
	api.RegisterTTSServer(routes, server.NewTTSServer(ttsService))

	// generator
	generatorService := service.NewGeneratorService(
		cfg,
		gameService,
		collectionService,
		deckService,
		cardService,
		ttsService,
		systemService,
	)
	api.RegisterGeneratorServer(routes, server.NewGeneratorServer(generatorService))

	// replace
	replaceService := service.NewReplaceService(ttsService)
	api.RegisterReplaceServer(routes, server.NewReplaceServer(replaceService))

	// recursive search
	searchService := service.NewSearchService(gameService, collectionService, deckService, cardService)
	api.RegisterSearchServer(routes, server.NewSearchServer(searchService))

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
