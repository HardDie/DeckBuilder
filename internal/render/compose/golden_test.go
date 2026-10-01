package compose_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/HardDie/DeckBuilder/internal/config"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/render/compose"
	oldgen "github.com/HardDie/DeckBuilder/internal/render/deprecated/generator"
	oldprogress "github.com/HardDie/DeckBuilder/internal/render/deprecated/progress"
	"github.com/HardDie/DeckBuilder/internal/render/generate/fake"
	"github.com/HardDie/DeckBuilder/internal/render/progress"
)

// Fixture is five 8×12 solid PNGs in testdata/fixture.
// One game "Raid", one collection "base", two decks.
// Bandits has one card (count 1), so the JSON emits a Card.
// Crew has Ada (count 2) and Bo (count 1), so the JSON emits a Deck.
// Back shadow is on. TTS card scale is 1.25, 1.5, 1.75. Request scale is 1.
// Deck order is by name, so Bandits precedes Crew.
//
// JSON goldens replace the results directory with RESULT and the clock
// with 2006-01-02 15:04:05. Image files are the raw bytes.
func TestIntegrationGoldenBytes(t *testing.T) {
	w := loadFixture(t)
	oldprogress.GetProgress().Flush()
	progress.Reset()

	oldDir := t.TempDir()
	run(t, oldDir, w, true)
	oldFiles := readResult(t, filepath.Join(oldDir, "result"))

	goldenDir := filepath.Join("testdata", "golden")
	if os.Getenv("WRITE_GOLDEN") == "1" {
		if err := os.MkdirAll(goldenDir, 0o755); err != nil {
			t.Fatal(err)
		}
		for name, body := range oldFiles {
			if err := os.WriteFile(filepath.Join(goldenDir, name), body, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}

	newDir := t.TempDir()
	run(t, newDir, w, false)
	newFiles := readResult(t, filepath.Join(newDir, "result"))

	want := readGolden(t, goldenDir)
	compare(t, "old", oldFiles, want)
	compare(t, "new", newFiles, want)
}

func loadFixture(t *testing.T) *fake.World {
	t.Helper()
	read := func(name string) []byte {
		t.Helper()
		body, err := os.ReadFile(filepath.Join("testdata", "fixture", name))
		if err != nil {
			t.Fatal(err)
		}
		return body
	}
	cfg := entitiesSettings.Default()
	cfg.EnableBackShadow = true
	cfg.CardSize.ScaleX = 1.25
	cfg.CardSize.ScaleY = 1.5
	cfg.CardSize.ScaleZ = 1.75
	return &fake.World{
		GameID:   "raid",
		GameName: "Raid",
		Settings: cfg,
		Collections: []*fake.Coll{{
			ID:   "base",
			Name: "Base",
			Decks: []*fake.Deck{
				{
					ID: "crew", Name: "Crew", Image: "crew.png", Back: read("crew_back.png"),
					Cards: []*fake.Card{
						{ID: 2, Name: "Ada", Description: "Pilot", Count: 2, Variables: map[string]string{"seat": "left"}, Face: read("crew_ada.png")},
						{ID: 3, Name: "Bo", Description: "Gunner", Count: 1, Face: read("crew_bo.png")},
					},
				},
				{
					ID: "bandits", Name: "Bandits", Image: "bandits.png", Back: read("bandits_back.png"),
					Cards: []*fake.Card{
						{ID: 1, Name: "Lookout", Description: "Watches the road", Count: 1, Variables: map[string]string{"role": "scout"}, Face: read("bandits_face.png")},
					},
				},
			},
		}},
	}
}

func run(t *testing.T, data string, w *fake.World, old bool) {
	t.Helper()
	cfg := config.Get("test")
	cfg.SetDataPath(data)
	reqSort := "name"
	var err error
	if old {
		err = oldgen.New(cfg, fake.Games{w}, fake.Collections{w}, fake.Decks{w}, fake.Cards{w}, fake.Systems{w}, fake.Speech{w}).
			GenerateGame("raid", oldgen.GenerateGameRequest{SortOrder: reqSort, Scale: 1})
	} else {
		err = compose.New(cfg, fake.Games{w}, fake.Collections{w}, fake.Decks{w}, fake.Cards{w}, fake.Systems{w}, fake.Speech{w}).
			GenerateGame("raid", compose.GenerateGameRequest{SortOrder: reqSort, Scale: 1})
	}
	if err != nil {
		t.Fatal(err)
	}
	if old {
		waitOld(t)
		return
	}
	waitDone(t)
}

func waitOld(t *testing.T) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		switch oldprogress.GetProgress().GetStatus().Status {
		case oldprogress.StatusDone:
			return
		case oldprogress.StatusError:
			t.Fatal(oldprogress.GetProgress().GetStatus().Message)
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("timed out waiting for the old generator")
}

func readResult(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	out := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		body, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if filepath.Ext(entry.Name()) == ".json" {
			body = normalizeJSON(body, abs)
		}
		out[entry.Name()] = body
	}
	return out
}

func readGolden(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	out := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		body, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		out[entry.Name()] = body
	}
	return out
}

var createdAt = regexp.MustCompile(`Created at: \d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)

func normalizeJSON(body []byte, results string) []byte {
	body = bytes.ReplaceAll(body, []byte(results), []byte("RESULT"))
	return createdAt.ReplaceAll(body, []byte("Created at: 2006-01-02 15:04:05"))
}

func compare(t *testing.T, label string, got, want map[string][]byte) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s has %d files, golden has %d (%v vs %v)", label, len(got), len(want), names(got), names(want))
	}
	for name, body := range want {
		g, ok := got[name]
		if !ok {
			t.Fatalf("%s missing %s", label, name)
		}
		if !bytes.Equal(g, body) {
			t.Fatalf("%s %s: %d bytes, golden %d bytes", label, name, len(g), len(body))
		}
	}
}

func names(files map[string][]byte) []string {
	out := make([]string, 0, len(files))
	for name := range files {
		out = append(out, name)
	}
	return out
}
