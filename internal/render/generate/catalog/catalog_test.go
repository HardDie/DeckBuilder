package catalog

import (
	"testing"

	"github.com/HardDie/DeckBuilder/internal/render/generate/fake"
)

func TestCollectSortsDecksAndSkipsEmpty(t *testing.T) {
	w := &fake.World{
		Collections: []*fake.Coll{{
			ID: "base",
			Decks: []*fake.Deck{
				{ID: "late", Name: "Zebra", Cards: []*fake.Card{{ID: 1, Count: 1}}},
				{ID: "empty", Name: "Middle"},
				{ID: "early", Name: "Apple", Cards: []*fake.Card{{ID: 2, Count: 3}}},
			},
		}},
	}
	grouped, order, err := Collect("raid", "name", fake.Collections{World: w}, fake.Decks{World: w}, fake.Cards{World: w})
	if err != nil {
		t.Fatal(err)
	}
	if len(order) != 2 || order[0].Name != "Apple" || order[1].Name != "Zebra" {
		t.Fatalf("order %+v", order)
	}
	if grouped[order[0]][0].Count != 3 || grouped[order[1]][0].ID != 1 {
		t.Fatalf("cards %+v", grouped)
	}
	sorts := w.Sorts()
	if len(sorts) != 5 {
		t.Fatalf("sort fields %v", sorts)
	}
	for _, field := range sorts {
		if field != "name" {
			t.Fatalf("sort fields %v", sorts)
		}
	}
}

func TestCollectMergesTheSameDeck(t *testing.T) {
	same := fake.Deck{ID: "crew", Name: "Crew", Image: "crew.png"}
	w := &fake.World{
		Collections: []*fake.Coll{
			{ID: "base", Decks: []*fake.Deck{{
				ID: same.ID, Name: same.Name, Image: same.Image,
				Cards: []*fake.Card{{ID: 1, Count: 1}},
			}}},
			{ID: "promo", Decks: []*fake.Deck{{
				ID: same.ID, Name: same.Name, Image: same.Image,
				Cards: []*fake.Card{{ID: 2, Count: 4}},
			}}},
		},
	}
	grouped, order, err := Collect("raid", "", fake.Collections{World: w}, fake.Decks{World: w}, fake.Cards{World: w})
	if err != nil {
		t.Fatal(err)
	}
	if len(order) != 1 {
		t.Fatalf("decks %d", len(order))
	}
	cards := grouped[order[0]]
	if len(cards) != 2 || cards[0].CollectionID != "base" || cards[1].CollectionID != "promo" {
		t.Fatalf("cards %+v", cards)
	}
}
