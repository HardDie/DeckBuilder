package script

import (
	"encoding/json"
	"testing"

	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
	entitiesSettings "github.com/HardDie/DeckBuilder/internal/entities/settings"
	"github.com/HardDie/DeckBuilder/internal/render/generate/catalog"
	"github.com/HardDie/DeckBuilder/internal/render/generate/fake"
	"github.com/HardDie/DeckBuilder/internal/tts_entity"
)

func TestBuildOneCardAndADeck(t *testing.T) {
	w := &fake.World{
		Collections: []*fake.Coll{{
			ID: "base",
			Decks: []*fake.Deck{
				{
					ID: "crew", Name: "Crew",
					Cards: []*fake.Card{
						{ID: 2, Name: "Ada", Description: "Pilot", Count: 2, Variables: map[string]string{"seat": "left"}},
						{ID: 3, Name: "Bo", Description: "Gunner", Count: 1},
					},
				},
				{
					ID: "bandits", Name: "Bandits",
					Cards: []*fake.Card{
						{ID: 1, Name: "Lookout", Description: "Watches the road", Count: 1, Variables: map[string]string{"role": "scout"}},
					},
				},
			},
		}},
	}
	grouped, order, err := catalog.Collect("raid", "name", fake.Collections{w}, fake.Decks{w}, fake.Cards{w})
	if err != nil {
		t.Fatal(err)
	}
	cfg := entitiesSettings.Default()
	cfg.CardSize.ScaleX = 1.25
	cfg.CardSize.ScaleY = 1.5
	cfg.CardSize.ScaleZ = 1.75
	pages := []Page{
		{DeckID: "bandits", Index: 1, Image: "/sheets/bandits.jpg", Back: "/sheets/bandits.png", Cols: 2, Rows: 2},
		{DeckID: "crew", Index: 1, Image: "/sheets/crew.jpg", Back: "/sheets/crew.png", Cols: 2, Rows: 2},
	}
	doc, err := Build(t.TempDir(), &entitiesGame.Game{ID: "raid", Name: "Raid"}, grouped, order, pages, &cfg, fake.Cards{w})
	if err != nil {
		t.Fatal(err)
	}
	if doc.Root.ObjectStates[0].Nickname != "Raid" {
		t.Fatalf("game %s", doc.Root.ObjectStates[0].Nickname)
	}
	col := doc.Root.ObjectStates[0].ContainedObjects[0].(*tts_entity.Bag)
	if col.Nickname != "base" || len(col.ContainedObjects) != 2 {
		t.Fatalf("collection %+v", col.Nickname)
	}
	card := col.ContainedObjects[0].(tts_entity.Card)
	if card.Name != "Card" || card.Nickname != "Lookout" || card.CardID != 100 || card.LuaScript != `role="scout"` {
		t.Fatalf("card %+v", card)
	}
	deck := col.ContainedObjects[1].(tts_entity.DeckObject)
	if deck.Name != "Deck" || deck.Nickname != "Crew" || len(deck.DeckIDs) != 3 {
		t.Fatalf("deck %+v ids %v", deck.Name, deck.DeckIDs)
	}
	if deck.DeckIDs[0] != 200 || deck.DeckIDs[2] != 201 || deck.Transform.ScaleX != 1.25 {
		t.Fatalf("ids %v scale %v", deck.DeckIDs, deck.Transform.ScaleX)
	}
	face := card.CustomDeck[1].FaceURL
	if face != "file:////sheets/bandits.jpg" {
		t.Fatalf("face %s", face)
	}
}

// The same catalog gives the same JSON on every build:
// collection bags by collection id, Lua lines by variable name.
func TestBuildIsStable(t *testing.T) {
	deck := func(id string) *fake.Deck {
		return &fake.Deck{ID: id, Name: id, Cards: []*fake.Card{
			{ID: 1, Name: "Ada", Count: 1, Variables: map[string]string{"hp": "2", "atk": "1", "def": "3", "name": "Ada"}},
		}}
	}
	w := &fake.World{Collections: []*fake.Coll{
		{ID: "zeta", Decks: []*fake.Deck{deck("z")}},
		{ID: "alpha", Decks: []*fake.Deck{deck("a")}},
		{ID: "mid", Decks: []*fake.Deck{deck("m")}},
	}}
	grouped, order, err := catalog.Collect("raid", "name", fake.Collections{World: w}, fake.Decks{World: w}, fake.Cards{World: w})
	if err != nil {
		t.Fatal(err)
	}
	var pages []Page
	for _, id := range []string{"z", "a", "m"} {
		pages = append(pages, Page{DeckID: id, Index: 1, Image: "/" + id + ".jpg", Back: "/" + id + ".png", Cols: 2, Rows: 2})
	}
	cfg := entitiesSettings.Default()

	build := func() (string, *tts_entity.Bag) {
		doc, err := Build(t.TempDir(), &entitiesGame.Game{ID: "raid", Name: "Raid"}, grouped, order, pages, &cfg, fake.Cards{World: w})
		if err != nil {
			t.Fatal(err)
		}
		doc.Root.ObjectStates[0].Description = "" // the clock
		body, err := json.Marshal(doc.Root)
		if err != nil {
			t.Fatal(err)
		}
		return string(body), &doc.Root.ObjectStates[0]
	}
	first, bag := build()
	for i := 0; i < 20; i++ {
		if got, _ := build(); got != first {
			t.Fatalf("build %d differs from the first", i+1)
		}
	}

	var nicknames []string
	for _, obj := range bag.ContainedObjects {
		nicknames = append(nicknames, obj.(*tts_entity.Bag).Nickname)
	}
	if want := []string{"alpha", "mid", "zeta"}; len(nicknames) != 3 || nicknames[0] != want[0] || nicknames[1] != want[1] || nicknames[2] != want[2] {
		t.Fatalf("collection bags %v, want %v", nicknames, want)
	}
	card := bag.ContainedObjects[0].(*tts_entity.Bag).ContainedObjects[0].(tts_entity.Card)
	if want := "atk=\"1\"\ndef=\"3\"\nhp=\"2\"\nname=\"Ada\""; card.LuaScript != want {
		t.Fatalf("LuaScript %q, want %q", card.LuaScript, want)
	}
}
