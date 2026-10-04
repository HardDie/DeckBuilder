//go:build !nomain

package main

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/render/compose"
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
	servicesReplace "github.com/HardDie/DeckBuilder/internal/services/replace"
	servicesSearch "github.com/HardDie/DeckBuilder/internal/services/search"
	servicesSystem "github.com/HardDie/DeckBuilder/internal/services/system"
	servicesTTS "github.com/HardDie/DeckBuilder/internal/services/tts"
	"github.com/HardDie/DeckBuilder/pkg/logger"
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
	generator  compose.Generator
	replace    servicesReplace.Replace
}

func wire(version string) *services {
	cfg := config.Get(version)

	// The log file comes first, so everything after it is recorded (ADR 026).
	logPath, err := logger.Init(logger.Options{
		Dir:     cfg.Logs(),
		File:    "deckbuilder.log",
		Format:  logger.FormatJSON,
		MaxSize: 5 << 20,
		Keep:    3,
	})
	if err != nil {
		slog.Warn("log file is off, logging to the console only", "err", err)
	}
	// One line per launch, so each run is easy to find in the log.
	// It goes before the log level setting is applied, so Warn or Error never hides it.
	slog.Info("app started", "version", version, "os", runtime.GOOS, "arch", runtime.GOARCH,
		"data", cfg.Data, "log", logPath)

	db := fsentry.New(cfg.Data, fsentry.WithPretty())
	if err := db.Init(); err != nil {
		slog.Error("open data folder", "err", err)
		os.Exit(1)
	}

	core := repositoriesCore.New(db)
	if err := core.Init(); err != nil {
		slog.Error("prepare catalog", "err", err)
		os.Exit(1)
	}

	repositorySettings := repositoriesSettings.New(cfg, db)
	serviceSystem := servicesSystem.New(repositorySettings)
	// The log level is a setting; until it is read, the log writes Info.
	if settings, err := serviceSystem.GetSettings(); err == nil {
		logger.SetLevel(settings.SlogLevel())
	} else {
		slog.Warn("read settings for the log level", "err", err)
	}

	repositoryGame := repositoriesGame.New(db)
	serviceGame := servicesGame.New(cfg, repositoryGame)

	repositoryCollection := repositoriesCollection.New(db)
	serviceCollection := servicesCollection.New(cfg, repositoryCollection)

	repositoryDeck := repositoriesDeck.New(db)
	serviceDeck := servicesDeck.New(cfg, repositoryDeck)

	repositoryCard := repositoriesCard.New(cfg, db)
	serviceCard := servicesCard.New(cfg, repositoryCard)

	serviceTTS := servicesTTS.New()
	serviceGenerator := compose.New(cfg, serviceGame, serviceCollection, serviceDeck, serviceCard, serviceSystem, serviceTTS)
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
