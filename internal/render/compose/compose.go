// Run a render: plan it, then draw the sheets.
package compose

import (
	"sync/atomic"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/render/generate"
	"github.com/HardDie/DeckBuilder/internal/render/generate/catalog"
	renderprogress "github.com/HardDie/DeckBuilder/internal/render/progress"
	"github.com/HardDie/DeckBuilder/internal/render/sheet/write"
	servicesCard "github.com/HardDie/DeckBuilder/internal/services/card"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
	servicesSystem "github.com/HardDie/DeckBuilder/internal/services/system"
	servicesTTS "github.com/HardDie/DeckBuilder/internal/services/tts"
	"github.com/HardDie/DeckBuilder/internal/tts_entity"
)

type Generator interface {
	GenerateGame(gameID string, req GenerateGameRequest) error
}

type GenerateGameRequest struct {
	SortOrder string
	Scale     int
}

type runner struct {
	cfg               *config.Config
	serviceGame       servicesGame.Game
	serviceCollection servicesCollection.Collection
	serviceDeck       servicesDeck.Deck
	serviceCard       servicesCard.Card
	serviceSystem     servicesSystem.System
	serviceTTS        servicesTTS.TTS

	// running is true from the start of GenerateGame until the run ends.
	// A second call meanwhile gets GenerateInProgress.
	running atomic.Bool
}

func New(
	cfg *config.Config,
	serviceGame servicesGame.Game,
	serviceCollection servicesCollection.Collection,
	serviceDeck servicesDeck.Deck,
	serviceCard servicesCard.Card,
	serviceSystem servicesSystem.System,
	serviceTTS servicesTTS.TTS,
) Generator {
	return &runner{
		cfg:               cfg,
		serviceGame:       serviceGame,
		serviceCollection: serviceCollection,
		serviceDeck:       serviceDeck,
		serviceCard:       serviceCard,
		serviceSystem:     serviceSystem,
		serviceTTS:        serviceTTS,
	}
}

func (s *runner) GenerateGame(gameID string, req GenerateGameRequest) error {
	if !s.running.CompareAndSwap(false, true) {
		return er.GenerateInProgress
	}
	// Release on an early error. Once the goroutine starts, it releases instead.
	started := false
	defer func() {
		if !started {
			s.running.Store(false)
		}
	}()

	cfg, err := s.serviceSystem.GetSettings()
	if err != nil {
		logger.Error.Printf("can't get config")
		return err
	}
	gameItem, err := s.serviceGame.Item(gameID)
	if err != nil {
		return err
	}
	decks, order, err := catalog.Collect(gameItem.ID, req.SortOrder, s.serviceCollection, s.serviceDeck, s.serviceCard)
	if err != nil {
		return err
	}
	if err := fs.RemoveFolder(s.cfg.Results()); err != nil {
		return err
	}
	if err := fs.CreateFolder(s.cfg.Results()); err != nil {
		return err
	}
	renderprogress.Begin()
	started = true
	go func() {
		defer s.running.Store(false)
		err := s.run(gameItem, decks, order, req.Scale, cfg)
		if err != nil {
			renderprogress.Fail()
			logger.Error.Println("Generator:", err.Error())
			return
		}
		renderprogress.Finish()
	}()
	return nil
}

func (s *runner) run(
	gameItem *entitiesGame.Game,
	decks map[catalog.Deck][]catalog.Card,
	order []catalog.Deck,
	scale int,
	cfg *entitiesSettings.Settings,
) error {
	plan, err := generate.Prepare(s.cfg.Results(), gameItem, decks, order, scale, cfg, s.serviceDeck, s.serviceCard)
	if err != nil {
		return err
	}
	total := len(plan.Sheets)
	renderprogress.Sheets(0, total)
	for _, back := range plan.Backs {
		if err := fs.CreateAndProcess(back.Path, back.Body, fs.BinToWriter); err != nil {
			return err
		}
	}
	for i, sheet := range plan.Sheets {
		if err := write.Draw(sheet.Faces, sheet.Back, sheet.CellW, sheet.CellH, sheet.Shadow, sheet.Path); err != nil {
			return err
		}
		renderprogress.Sheets(i+1, total)
	}
	if err := fs.CreateAndProcess(plan.JSONPath, plan.Root, fs.JsonToWriter[tts_entity.RootObjects]); err != nil {
		return err
	}
	s.serviceTTS.SendToTTS(plan.Bag)
	return nil
}
