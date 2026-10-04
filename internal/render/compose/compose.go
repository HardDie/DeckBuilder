// Run a render: plan it, then draw the sheets.
package compose

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/fs"
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
		return apperr.ErrRenderInProgress
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
		slog.Error("read settings for render", "err", err)
		return err
	}
	gameItem, err := s.serviceGame.Item(gameID)
	if err != nil {
		return err
	}
	start := time.Now()
	decks, order, err := catalog.Collect(gameItem.ID, req.SortOrder, s.serviceCollection, s.serviceDeck, s.serviceCard)
	if err != nil {
		return err
	}
	collected := time.Now()
	// Every deck needs a back and every card a face. Say what is missing instead of starting.
	missing, err := generate.MissingImages(decks, order, s.serviceDeck)
	if err != nil {
		return err
	}
	slog.Debug("render catalog", "game", gameItem.ID, "decks", len(order), "cards", countCards(decks),
		"collect_ms", collected.Sub(start).Milliseconds(),
		"check_images_ms", time.Since(collected).Milliseconds())
	if len(missing) > 0 {
		return apperr.With(apperr.ErrMissingImages, missingMessage(missing))
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
		err := s.safeRun(dir, gameItem, decks, order, req.Scale, cfg, start)
		if err != nil {
			renderprogress.Fail(err)
			slog.Error("render failed", "game", gameItem.ID, "err", err)
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
	start time.Time,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v\n%s", r, debug.Stack())
		}
	}()
	return s.run(dir, gameItem, decks, order, scale, cfg, start)
}

// countCards is the number of catalog cards over all decks.
func countCards(decks map[catalog.Deck][]catalog.Card) int {
	n := 0
	for _, cards := range decks {
		n += len(cards)
	}
	return n
}

// run draws the plan into dir.
// A file whose name is already in dir is reused: names carry a hash of their content.
// New files go through a temporary name, so a failed run never leaves a broken file.
// Files the plan does not use are removed only after everything was written.
// start is when GenerateGame began; each step is logged at Debug with its time.
func (s *runner) run(
	dir string,
	gameItem *entitiesGame.Game,
	decks map[catalog.Deck][]catalog.Card,
	order []catalog.Deck,
	scale int,
	cfg *entitiesSettings.Settings,
	start time.Time,
) error {
	step := time.Now()
	plan, err := generate.Prepare(dir, gameItem, decks, order, scale, cfg, s.serviceDeck, s.serviceCard)
	if err != nil {
		return err
	}
	prepareTime := time.Since(step)
	step = time.Now()
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
	backsTime := time.Since(step)
	step = time.Now()
	reused := 0
	var drawn write.Timings
	for i, sheet := range plan.Sheets {
		keep[filepath.Base(sheet.Path)] = struct{}{}
		if fs.FileExists(sheet.Path) {
			reused++
			slog.Debug("render sheet reused", "deck", sheet.DeckName, "file", filepath.Base(sheet.Path))
		} else {
			sheetStart := time.Now()
			var t write.Timings
			err := fs.WriteAtomic(sheet.Path, func(tmp string) error {
				var err error
				t, err = write.Draw(sheet.Faces, sheet.Back, sheet.CellW, sheet.CellH, sheet.Shadow, tmp)
				return err
			})
			slog.Debug("render sheet", "deck", sheet.DeckName, "file", filepath.Base(sheet.Path),
				"faces", len(sheet.Faces), "total_ms", time.Since(sheetStart).Milliseconds(),
				"decode_resize_ms", t.Jobs.Milliseconds(),
				"decode_sum_ms", t.Decode.Milliseconds(), "resize_sum_ms", t.Resize.Milliseconds(),
				"paint_ms", t.Paint.Milliseconds(), "encode_ms", t.Encode.Milliseconds(),
				"write_ms", t.Write.Milliseconds())
			drawn = addTimings(drawn, t)
			if errors.Is(err, write.ErrUndecodable) {
				// The sheet knows its deck, not the card; the decoder's detail goes to the log.
				slog.Warn("sheet image does not decode", "deck", sheet.DeckName, "err", err)
				return generate.UnreadableImage(sheet.DeckName)
			}
			if err != nil {
				return err
			}
		}
		renderprogress.Sheets(i+1, total)
	}
	sheetsTime := time.Since(step)
	step = time.Now()
	keep[filepath.Base(plan.JSONPath)] = struct{}{}
	err = fs.WriteAtomic(plan.JSONPath, func(tmp string) error {
		return fs.CreateAndProcess(tmp, plan.Root, fs.JsonToWriter[tts_entity.RootObjects])
	})
	if err != nil {
		return err
	}
	removeStale(dir, keep)
	filesTime := time.Since(step)
	step = time.Now()
	s.serviceTTS.SendToTTS(plan.Bag)
	slog.Debug("render steps", "game", gameItem.ID,
		"total_ms", time.Since(start).Milliseconds(), "prepare_ms", prepareTime.Milliseconds(), "backs_ms", backsTime.Milliseconds(),
		"sheets_ms", sheetsTime.Milliseconds(),
		"decode_resize_ms", drawn.Jobs.Milliseconds(),
		"decode_sum_ms", drawn.Decode.Milliseconds(), "resize_sum_ms", drawn.Resize.Milliseconds(),
		"paint_ms", drawn.Paint.Milliseconds(), "encode_ms", drawn.Encode.Milliseconds(),
		"write_ms", drawn.Write.Milliseconds(),
		"json_and_cleanup_ms", filesTime.Milliseconds(), "tts_ms", time.Since(step).Milliseconds())
	slog.Info("render finished", "game", gameItem.ID, "sheets", total, "reused", reused,
		"duration_ms", time.Since(start).Milliseconds())
	return nil
}

// addTimings sums two sheets' timings.
func addTimings(a, b write.Timings) write.Timings {
	return write.Timings{
		Jobs:   a.Jobs + b.Jobs,
		Decode: a.Decode + b.Decode,
		Resize: a.Resize + b.Resize,
		Paint:  a.Paint + b.Paint,
		Encode: a.Encode + b.Encode,
		Write:  a.Write + b.Write,
	}
}

// removeStale deletes the files in dir that the last render did not write or reuse.
// It is best effort: a file that cannot be removed is logged and left.
func removeStale(dir string, keep map[string]struct{}) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		slog.Warn("list stale render files", "dir", dir, "err", err)
		return
	}
	for _, e := range entries {
		if _, ok := keep[e.Name()]; ok || e.IsDir() {
			continue
		}
		if err := os.Remove(filepath.Join(dir, e.Name())); err != nil {
			slog.Warn("remove stale render file", "file", e.Name(), "err", err)
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
			slog.Warn("remove old result file", "file", e.Name(), "err", err)
		}
	}
}
