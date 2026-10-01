package compose_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/render/compose"
	"github.com/HardDie/DeckBuilder/internal/render/generate/fake"
	"github.com/HardDie/DeckBuilder/internal/render/progress"
)

func TestGenerateUnknownGame(t *testing.T) {
	progress.Reset()
	dir := t.TempDir()
	cfg := config.Get("test")
	cfg.SetDataPath(dir)
	w := &fake.World{GameID: "raid", GameName: "Raid"}
	err := compose.New(cfg, fake.Games{w}, fake.Collections{w}, fake.Decks{w}, fake.Cards{w}, fake.Systems{w}, fake.Speech{w}).
		GenerateGame("missing", compose.GenerateGameRequest{})
	if err == nil {
		t.Fatal("expected missing game")
	}
	if progress.Get().Status != progress.Empty {
		t.Fatalf("status %s", progress.Get().Status)
	}
	if _, err := os.Stat(cfg.Results()); !os.IsNotExist(err) {
		t.Fatal("results folder was created")
	}
}

func TestGenerateClearsResults(t *testing.T) {
	progress.Reset()
	dir := t.TempDir()
	cfg := config.Get("test")
	cfg.SetDataPath(dir)
	if err := os.MkdirAll(cfg.Results(), 0o755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(cfg.Results(), "old.txt")
	if err := os.WriteFile(marker, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	w := raidWorld()
	err := compose.New(cfg, fake.Games{w}, fake.Collections{w}, fake.Decks{w}, fake.Cards{w}, fake.Systems{w}, fake.Speech{w}).
		GenerateGame("raid", compose.GenerateGameRequest{SortOrder: "name", Scale: 1})
	if err != nil {
		t.Fatal(err)
	}
	waitDone(t)
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("old result file stayed")
	}
	if _, err := os.Stat(filepath.Join(cfg.Results(), "raid.json")); err != nil {
		t.Fatal(err)
	}
	if len(w.TTS()) != 1 {
		t.Fatalf("tts calls %d", len(w.TTS()))
	}
}

func TestReportOneSheet(t *testing.T) {
	runReport(t, raidWorld())
	got := progress.Get()
	if got.Done != 1 || got.Total != 1 || got.Percent != 100 || got.Status != progress.Done {
		t.Fatalf("one sheet %+v", got)
	}
}

func TestReportManySheets(t *testing.T) {
	w := raidWorld()
	deck := w.Collections[0].Decks[0]
	w.Collections[0].Decks = append(w.Collections[0].Decks, &fake.Deck{
		ID: "crew", Name: "Crew", Back: deck.Back,
		Cards: []*fake.Card{{ID: 2, Name: "Ada", Count: 1, Face: deck.Cards[0].Face}},
	})
	runReport(t, w)
	got := progress.Get()
	if got.Done != 2 || got.Total != 2 || got.Percent != 100 || got.Status != progress.Done {
		t.Fatalf("many sheets %+v", got)
	}
}

func runReport(t *testing.T, w *fake.World) {
	t.Helper()
	progress.Reset()
	dir := t.TempDir()
	cfg := config.Get("test")
	cfg.SetDataPath(dir)
	err := compose.New(cfg, fake.Games{w}, fake.Collections{w}, fake.Decks{w}, fake.Cards{w}, fake.Systems{w}, fake.Speech{w}).
		GenerateGame("raid", compose.GenerateGameRequest{SortOrder: "name", Scale: 1})
	if err != nil {
		t.Fatal(err)
	}
	waitDone(t)
}

func raidWorld() *fake.World {
	face := solidPNG(8, 12, color.RGBA{R: 180, A: 255})
	back := solidPNG(8, 12, color.RGBA{B: 180, A: 255})
	cfg := entitiesSettings.Default()
	return &fake.World{
		GameID:   "raid",
		GameName: "Raid",
		Settings: cfg,
		Collections: []*fake.Coll{{
			ID: "base",
			Decks: []*fake.Deck{{
				ID: "bandits", Name: "Bandits", Back: back,
				Cards: []*fake.Card{{ID: 1, Name: "Lookout", Count: 1, Face: face}},
			}},
		}},
	}
}

func solidPNG(w, h int, c color.RGBA) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func waitDone(t *testing.T) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		switch progress.Get().Status {
		case progress.Done:
			return
		case progress.Error:
			t.Fatal("generate failed")
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("timed out waiting for generate")
}
