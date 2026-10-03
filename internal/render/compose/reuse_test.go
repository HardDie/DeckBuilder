package compose_test

import (
	"bytes"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/render/compose"
	"github.com/HardDie/DeckBuilder/internal/render/generate/fake"
	"github.com/HardDie/DeckBuilder/internal/render/progress"
)

var oldTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

// twoDecks is one game with two one-page decks, so each can change alone.
func twoDecks() *fake.World {
	back := solidPNG(8, 12, color.RGBA{B: 180, A: 255})
	return &fake.World{
		GameID:   "raid",
		GameName: "Raid",
		Settings: entitiesSettings.Default(),
		Collections: []*fake.Coll{{
			ID: "base",
			Decks: []*fake.Deck{
				{ID: "bandits", Name: "Bandits", Back: back, Cards: []*fake.Card{
					{ID: 1, Name: "Lookout", Count: 1, Face: solidPNG(8, 12, color.RGBA{R: 180, A: 255})},
				}},
				{ID: "crew", Name: "Crew", Back: back, Cards: []*fake.Card{
					{ID: 2, Name: "Ada", Count: 1, Face: solidPNG(8, 12, color.RGBA{G: 180, A: 255})},
					{ID: 3, Name: "Bo", Count: 1, Face: solidPNG(8, 12, color.RGBA{R: 90, G: 90, A: 255})},
				}},
			},
		}},
	}
}

func generate(t *testing.T, cfg *config.Config, w *fake.World) string {
	t.Helper()
	err := compose.New(cfg, fake.Games{World: w}, fake.Collections{World: w}, fake.Decks{World: w}, fake.Cards{World: w}, fake.Systems{World: w}, fake.Speech{World: w}).
		GenerateGame("raid", compose.GenerateGameRequest{SortOrder: "name", Scale: 1})
	if err != nil {
		t.Fatal(err)
	}
	return waitFinished(t)
}

// files lists the game folder: name → modification time.
func files(t *testing.T, dir string) map[string]time.Time {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	out := make(map[string]time.Time, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			t.Fatal(err)
		}
		out[e.Name()] = info.ModTime()
	}
	return out
}

// backdate sets every file to oldTime, so a reused file keeps that time.
func backdate(t *testing.T, dir string) {
	t.Helper()
	for name := range files(t, dir) {
		if err := os.Chtimes(filepath.Join(dir, name), oldTime, oldTime); err != nil {
			t.Fatal(err)
		}
	}
}

func withPrefix(m map[string]time.Time, prefix string) []string {
	var out []string
	for name := range m {
		if strings.HasPrefix(name, prefix) {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

func newCfg(t *testing.T) (*config.Config, string) {
	t.Helper()
	progress.Reset()
	cfg := config.Get("test")
	cfg.SetDataPath(t.TempDir())
	return cfg, filepath.Join(cfg.Results(), "raid")
}

func TestGenerateReusesUnchangedPages(t *testing.T) {
	cfg, dir := newCfg(t)
	w := twoDecks()
	if status := generate(t, cfg, w); status != progress.Done {
		t.Fatalf("first run %s", status)
	}
	first := files(t, dir)
	if len(withPrefix(first, "bandits_")) != 1 || len(withPrefix(first, "crew_")) != 1 || len(withPrefix(first, "backside_")) != 2 {
		t.Fatalf("first run files %v", first)
	}

	t.Run("nothing_changed", func(t *testing.T) {
		backdate(t, dir)
		generate(t, cfg, w)
		got := files(t, dir)
		for name, mod := range got {
			if name == "raid.json" {
				continue
			}
			if !mod.Equal(oldTime) {
				t.Errorf("%s was redrawn", name)
			}
		}
		if len(got) != len(first) {
			t.Fatalf("files %v, want %v", got, first)
		}
	})

	t.Run("card_text_is_not_drawn", func(t *testing.T) {
		backdate(t, dir)
		before, err := os.ReadFile(filepath.Join(dir, "raid.json"))
		if err != nil {
			t.Fatal(err)
		}
		ada := w.Collections[0].Decks[1].Cards[0]
		ada.Name, ada.Count, ada.Variables = "Ada Prime", 3, map[string]string{"hp": "2"}

		generate(t, cfg, w)
		got := files(t, dir)
		for _, name := range append(withPrefix(got, "crew_"), withPrefix(got, "backside_")...) {
			if !got[name].Equal(oldTime) {
				t.Errorf("%s was redrawn for a text change", name)
			}
		}
		after, err := os.ReadFile(filepath.Join(dir, "raid.json"))
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(before, after) {
			t.Fatal("the JSON did not pick up the new card text")
		}
	})

	t.Run("one_face_changed", func(t *testing.T) {
		backdate(t, dir)
		before := files(t, dir)
		w.Collections[0].Decks[1].Cards[1].Face = solidPNG(8, 12, color.RGBA{G: 40, B: 200, A: 255})

		generate(t, cfg, w)
		got := files(t, dir)
		if b := withPrefix(got, "bandits_"); len(b) != 1 || !got[b[0]].Equal(oldTime) {
			t.Errorf("bandits page should be reused: %v", b)
		}
		crew := withPrefix(got, "crew_")
		if len(crew) != 1 || crew[0] == withPrefix(before, "crew_")[0] {
			t.Fatalf("crew page should get a new name: %v (before %v)", crew, withPrefix(before, "crew_"))
		}
		if len(got) != len(before) {
			t.Fatalf("the old crew page was not removed: %v", got)
		}
	})
}

func TestGenerateFailureKeepsPreviousResult(t *testing.T) {
	cfg, dir := newCfg(t)
	w := twoDecks()
	if status := generate(t, cfg, w); status != progress.Done {
		t.Fatalf("first run %s", status)
	}
	before := files(t, dir)

	// A new face that cannot be decoded: its page must be drawn, and drawing fails.
	w.Collections[0].Decks[0].Cards[0].Face = []byte("not an image")
	if status := generate(t, cfg, w); status != progress.Error {
		t.Fatalf("second run %s, want error", status)
	}
	got := files(t, dir)
	if len(got) != len(before) {
		t.Fatalf("files %v, want the previous %v", got, before)
	}
	for name := range before {
		if _, ok := got[name]; !ok {
			t.Errorf("%s was removed by a failed run", name)
		}
	}
	for name := range got {
		if strings.HasSuffix(name, ".tmp") {
			t.Errorf("temporary file %s left behind", name)
		}
	}
}
