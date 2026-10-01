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
	"github.com/HardDie/DeckBuilder/internal/progress"
	"github.com/HardDie/DeckBuilder/internal/render/compose"
	"github.com/HardDie/DeckBuilder/internal/render/generate/fake"
)

func TestGenerateUnknownGame(t *testing.T) {
	progress.GetProgress().Flush()
	dir := t.TempDir()
	cfg := config.Get("test")
	cfg.SetDataPath(dir)
	w := &fake.World{GameID: "raid", GameName: "Raid"}
	err := compose.New(cfg, fake.Games{w}, fake.Collections{w}, fake.Decks{w}, fake.Cards{w}, fake.Systems{w}, fake.Speech{w}).
		GenerateGame("missing", compose.GenerateGameRequest{})
	if err == nil {
		t.Fatal("expected missing game")
	}
	if progress.GetProgress().GetStatus().Status != progress.StatusEmpty {
		t.Fatalf("status %s", progress.GetProgress().GetStatus().Status)
	}
	if _, err := os.Stat(cfg.Results()); !os.IsNotExist(err) {
		t.Fatal("results folder was created")
	}
}

func TestGenerateClearsResults(t *testing.T) {
	progress.GetProgress().Flush()
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
		switch progress.GetProgress().GetStatus().Status {
		case progress.StatusDone:
			return
		case progress.StatusError:
			t.Fatal(progress.GetProgress().GetStatus().Message)
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("timed out waiting for generate")
}
