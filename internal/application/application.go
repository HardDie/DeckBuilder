package application

import (
	"net/http"

	"github.com/HardDie/fsentry"
	"github.com/gorilla/mux"

	"github.com/HardDie/DeckBuilder/internal/api"
	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/db"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/repository"
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
	builderDB := db.NewFSEntryDB(fsentry.NewFSEntry(cfg.Games()))

	// game
	gameRepository := repository.NewGameRepository(cfg, builderDB)
	gameService := service.NewGameService(cfg, gameRepository)
	api.RegisterGameServer(routes, server.NewGameServer(gameService))

	// collection
	collectionRepository := repository.NewCollectionRepository(cfg, builderDB)
	collectionService := service.NewCollectionService(cfg, collectionRepository)
	api.RegisterCollectionServer(routes, server.NewCollectionServer(collectionService))

	// deck
	deckRepository := repository.NewDeckRepository(cfg, collectionRepository)
	deckService := service.NewDeckService(deckRepository)
	api.RegisterDeckServer(routes, server.NewDeckServer(deckService))

	// card
	cardService := service.NewCardService(repository.NewCardRepository(cfg, deckRepository))
	api.RegisterCardServer(routes, server.NewCardServer(cardService))

	// image
	api.RegisterImageServer(routes, server.NewImageServer(gameService, collectionService, deckService, cardService))

	// system
	api.RegisterSystemServer(routes, server.NewSystemServer(cfg))

	// generator
	generatorService := service.NewGeneratorService(cfg, gameService, collectionService, deckService, cardService)
	api.RegisterGeneratorServer(routes, server.NewGeneratorServer(generatorService))

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
