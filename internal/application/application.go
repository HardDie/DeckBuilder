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
	builderDB := db.NewFSEntryDB(fsentry.NewFSEntry(cfg.Data))
	err := builderDB.Init()
	if err != nil {
		logger.Error.Fatal(err)
	}

	// system
	systemServer := server.NewSystemServer(cfg, service.NewService(cfg, builderDB))
	api.RegisterSystemServer(routes, systemServer)

	// game
	gameRepository := repository.NewGameRepository(cfg, builderDB)
	gameService := service.NewGameService(cfg, gameRepository)
	api.RegisterGameServer(routes, server.NewGameServer(gameService, systemServer))

	// collection
	collectionRepository := repository.NewCollectionRepository(cfg, builderDB)
	collectionService := service.NewCollectionService(cfg, collectionRepository)
	api.RegisterCollectionServer(routes, server.NewCollectionServer(collectionService, systemServer))

	// deck
	deckRepository := repository.NewDeckRepository(cfg, builderDB)
	deckService := service.NewDeckService(cfg, deckRepository)
	api.RegisterDeckServer(routes, server.NewDeckServer(deckService, systemServer))

	// card
	cardService := service.NewCardService(cfg, repository.NewCardRepository(cfg, builderDB))
	api.RegisterCardServer(routes, server.NewCardServer(cardService, systemServer))

	// image
	api.RegisterImageServer(routes, server.NewImageServer(gameService, collectionService, deckService, cardService))

	// generator
	generatorService := service.NewGeneratorService(cfg, gameService, collectionService, deckService, cardService)
	api.RegisterGeneratorServer(routes, server.NewGeneratorServer(generatorService))

	// replace
	replaceService := service.NewReplaceService()
	api.RegisterReplaceServer(routes, server.NewReplaceServer(replaceService))

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
