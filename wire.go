//go:build !nomain

package main

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

type services struct {
	cfg        *config.Config
	game       servicesGame.Game
	collection servicesCollection.Collection
	deck       servicesDeck.Deck
	card       servicesCard.Card
	system     servicesSystem.System
	tts        servicesTTS.TTS
	search     servicesSearch.Search
	generator  servicesGenerator.Generator
	replace    servicesReplace.Replace
}

func wire(version string) *services {
	cfg := config.Get(version)

	db := fsentry.New(cfg.Data, fsentry.WithPretty())
	if err := db.Init(); err != nil {
		logger.Error.Fatal(err)
	}

	core := repositoriesCore.New(db)
	if err := core.Init(); err != nil {
		logger.Error.Fatal(err)
	}

	repositorySettings := repositoriesSettings.New(cfg, db)
	serviceSystem := servicesSystem.New(repositorySettings)

	repositoryGame := repositoriesGame.New(cfg, db)
	serviceGame := servicesGame.New(cfg, repositoryGame)

	repositoryCollection := repositoriesCollection.New(cfg, db)
	serviceCollection := servicesCollection.New(cfg, repositoryCollection)

	repositoryDeck := repositoriesDeck.New(cfg, db)
	serviceDeck := servicesDeck.New(cfg, repositoryDeck)

	repositoryCard := repositoriesCard.New(cfg, db)
	serviceCard := servicesCard.New(cfg, repositoryCard)

	serviceTTS := servicesTTS.New()
	serviceGenerator := servicesGenerator.New(cfg, serviceGame, serviceCollection, serviceDeck, serviceCard, serviceSystem, serviceTTS)
	serviceReplace := servicesReplace.New(serviceTTS)
	serviceSearch := servicesSearch.New(serviceGame, serviceCollection, serviceDeck, serviceCard)

	return &services{
		cfg:        cfg,
		game:       serviceGame,
		collection: serviceCollection,
		deck:       serviceDeck,
		card:       serviceCard,
		system:     serviceSystem,
		tts:        serviceTTS,
		search:     serviceSearch,
		generator:  serviceGenerator,
		replace:    serviceReplace,
	}
}
