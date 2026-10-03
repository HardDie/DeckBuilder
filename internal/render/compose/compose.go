// Run a render: plan it, then draw the sheets.
package compose

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
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
	// Every deck needs a back and every card a face. Say what is missing instead of starting.
	missing, err := generate.MissingImages(decks, order, s.serviceDeck)
	if err != nil {
		return err
	}
	if len(missing) > 0 {
		return er.GenerateMissingImages.AddMessage(missingMessage(missing))
	}
	// Each game renders into its own folder; the previous files stay until the run succeeds.
	dir := filepath.Join(s.cfg.Results(), gameItem.ID)
	removeLegacyFiles(s.cfg.Results())
	if err := fs.CreateFolder(dir); err != nil {
		return err
	}
	renderprogress.Begin()
	started = true
	go func() {
		defer s.running.Store(false)
		err := s.safeRun(dir, gameItem, decks, order, req.Scale, cfg)
		if err != nil {
			renderprogress.Fail()
			logger.Error.Println("Generator:", err.Error())
			return
		}
		renderprogress.Finish()
	}()
	return nil
}

// missingShown caps the list in the message; the rest is counted.
const missingShown = 5

// missingMessage is the text the window shows when images are missing.
func missingMessage(missing []string) string {
	shown := missing[:min(len(missing), missingShown)]
	msg := "Render needs an image for every deck and card. Missing: " + strings.Join(shown, "; ")
	if rest := len(missing) - len(shown); rest > 0 {
		msg += fmt.Sprintf("; and %d more", rest)
	}
	return msg + "."
}

// safeRun calls run and turns a panic into an error.
// The goroutine then fails the progress like any other error.
func (s *runner) safeRun(
	dir string,
	gameItem *entitiesGame.Game,
	decks map[catalog.Deck][]catalog.Card,
	order []catalog.Deck,
	scale int,
	cfg *entitiesSettings.Settings,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v\n%s", r, debug.Stack())
		}
	}()
	return s.run(dir, gameItem, decks, order, scale, cfg)
}

// run draws the plan into dir.
// A file whose name is already in dir is reused: names carry a hash of their content.
// New files go through a temporary name, so a failed run never leaves a broken file.
// Files the plan does not use are removed only after everything was written.
func (s *runner) run(
	dir string,
	gameItem *entitiesGame.Game,
	decks map[catalog.Deck][]catalog.Card,
	order []catalog.Deck,
	scale int,
	cfg *entitiesSettings.Settings,
) error {
	plan, err := generate.Prepare(dir, gameItem, decks, order, scale, cfg, s.serviceDeck, s.serviceCard)
	if err != nil {
		return err
	}
	keep := make(map[string]struct{})
	total := len(plan.Sheets)
	renderprogress.Sheets(0, total)
	for _, back := range plan.Backs {
		keep[filepath.Base(back.Path)] = struct{}{}
		if fs.FileExists(back.Path) {
			continue
		}
		err := fs.WriteAtomic(back.Path, func(tmp string) error {
			return fs.CreateAndProcess(tmp, back.Body, fs.BinToWriter)
		})
		if err != nil {
			return err
		}
	}
	reused := 0
	for i, sheet := range plan.Sheets {
		keep[filepath.Base(sheet.Path)] = struct{}{}
		if fs.FileExists(sheet.Path) {
			reused++
		} else {
			err := fs.WriteAtomic(sheet.Path, func(tmp string) error {
				return write.Draw(sheet.Faces, sheet.Back, sheet.CellW, sheet.CellH, sheet.Shadow, tmp)
			})
			if err != nil {
				return err
			}
		}
		renderprogress.Sheets(i+1, total)
	}
	keep[filepath.Base(plan.JSONPath)] = struct{}{}
	err = fs.WriteAtomic(plan.JSONPath, func(tmp string) error {
		return fs.CreateAndProcess(tmp, plan.Root, fs.JsonToWriter[tts_entity.RootObjects])
	})
	if err != nil {
		return err
	}
	removeStale(dir, keep)
	logger.Info.Printf("Generator: %d of %d sheets reused", reused, total)
	s.serviceTTS.SendToTTS(plan.Bag)
	return nil
}

// removeStale deletes the files in dir that the last render did not write or reuse.
// It is best effort: a file that cannot be removed is logged and left.
func removeStale(dir string, keep map[string]struct{}) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		logger.Warn.Println("Generator: list stale files:", err.Error())
		return
	}
	for _, e := range entries {
		if _, ok := keep[e.Name()]; ok || e.IsDir() {
			continue
		}
		if err := os.Remove(filepath.Join(dir, e.Name())); err != nil {
			logger.Warn.Println("Generator: remove stale file:", err.Error())
		}
	}
}

// removeLegacyFiles deletes files directly in result/.
// Before per-game folders, every render wrote there and cleared it first.
func removeLegacyFiles(results string) {
	entries, err := os.ReadDir(results)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if err := os.Remove(filepath.Join(results, e.Name())); err != nil {
			logger.Warn.Println("Generator: remove old result file:", err.Error())
		}
	}
}
