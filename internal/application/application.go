package application

import (
	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/logger"
	repositoriesCard "github.com/HardDie/DeckBuilder/internal/repositories/card"
	repositoriesCollection "github.com/HardDie/DeckBuilder/internal/repositories/collection"
	repositoriesCore "github.com/HardDie/DeckBuilder/internal/repositories/core"
	repositoriesDeck "github.com/HardDie/DeckBuilder/internal/repositories/deck"
	repositoriesGame "github.com/HardDie/DeckBuilder/internal/repositories/game"
	repositoriesSettings "github.com/HardDie/DeckBuilder/internal/repositories/settings"
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
	cfg               *config.Config
	serviceGame       servicesGame.Game
	serviceCollection servicesCollection.Collection
	serviceDeck       servicesDeck.Deck
	serviceCard       servicesCard.Card
	serviceSystem     servicesSystem.System
	serviceTTS        servicesTTS.TTS
	serviceSearch     servicesSearch.Search
	serviceGenerator  servicesGenerator.Generator
	serviceReplace    servicesReplace.Replace
}

func Get(version string) (*Application, error) {
	cfg := config.Get(version)

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

	// game
	repositoryGame := repositoriesGame.New(cfg, db)
	serviceGame := servicesGame.New(cfg, repositoryGame)

	// collection
	repositoryCollection := repositoriesCollection.New(cfg, db)
	serviceCollection := servicesCollection.New(cfg, repositoryCollection)

	// deck
	repositoryDeck := repositoriesDeck.New(cfg, db)
	serviceDeck := servicesDeck.New(cfg, repositoryDeck)

	// card
	repositoryCard := repositoriesCard.New(cfg, db)
	serviceCard := servicesCard.New(cfg, repositoryCard)

	serviceTTS := servicesTTS.New()

	// generator
	serviceGenerator := servicesGenerator.New(cfg, serviceGame, serviceCollection, serviceDeck, serviceCard, serviceSystem, serviceTTS)

	// replace
	serviceReplace := servicesReplace.New(serviceTTS)

	// recursive search
	serviceSearch := servicesSearch.New(serviceGame, serviceCollection, serviceDeck, serviceCard)

	return &Application{
		cfg:               cfg,
		serviceGame:       serviceGame,
		serviceCollection: serviceCollection,
		serviceDeck:       serviceDeck,
		serviceCard:       serviceCard,
		serviceSystem:     serviceSystem,
		serviceTTS:        serviceTTS,
		serviceSearch:     serviceSearch,
		serviceGenerator:  serviceGenerator,
		serviceReplace:    serviceReplace,
	}, nil
}

func (app *Application) Config() *config.Config {
	return app.cfg
}

func (app *Application) GameService() servicesGame.Game {
	return app.serviceGame
}

func (app *Application) CollectionService() servicesCollection.Collection {
	return app.serviceCollection
}

func (app *Application) DeckService() servicesDeck.Deck {
	return app.serviceDeck
}

func (app *Application) CardService() servicesCard.Card {
	return app.serviceCard
}

func (app *Application) SystemService() servicesSystem.System {
	return app.serviceSystem
}

func (app *Application) TTSService() servicesTTS.TTS {
	return app.serviceTTS
}

func (app *Application) SearchService() servicesSearch.Search {
	return app.serviceSearch
}

func (app *Application) GeneratorService() servicesGenerator.Generator {
	return app.serviceGenerator
}

func (app *Application) ReplaceService() servicesReplace.Replace {
	return app.serviceReplace
}
