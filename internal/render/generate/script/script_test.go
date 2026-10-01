package script

import (
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
