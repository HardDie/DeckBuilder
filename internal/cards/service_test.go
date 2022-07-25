package cards

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/decks"
	er "tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
)

var (
	gameID       = "test_game"
	collectionID = "test_collection"
	deckID       = "test_deck"
)

func testCreate(t *testing.T) {
	service := NewService()
	cardTitle := "one"
	desc := "best card ever"

	// Create card
	card, err := service.Create(gameID, collectionID, deckID, &CreateCardDTO{
		Title:       cardTitle,
		Description: desc,
	})
	if err != nil {
		t.Fatal(err)
	}
	if card.Title.String() != cardTitle {
		t.Fatal("Bad title [got]", card.Title, "[want]", cardTitle)
	}
	if card.Description.String() != desc {
		t.Fatal("Bad description [got]", card.Description, "[want]", desc)
	}

	// Delete card
	err = service.Delete(gameID, collectionID, deckID, card.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func testDelete(t *testing.T) {
	service := NewService()
	cardTitle := "one"

	// Try to remove non-existing card
	err := service.Delete(gameID, collectionID, deckID, 1)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card, err := service.Create(gameID, collectionID, deckID, &CreateCardDTO{
		Title: cardTitle,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete card
	err = service.Delete(gameID, collectionID, deckID, card.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete card twice
	err = service.Delete(gameID, collectionID, deckID, card.ID)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}
}
func testUpdate(t *testing.T) {
	service := NewService()
	cardTitle := []string{"one", "two"}
	desc := []string{"first description", "second description"}

	// Try to update non-existing card
	_, err := service.Update(gameID, collectionID, deckID, 1, &UpdateCardDTO{})
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card1, err := service.Create(gameID, collectionID, deckID, &CreateCardDTO{
		Title:       cardTitle[0],
		Description: desc[0],
	})
	if err != nil {
		t.Fatal(err)
	}
	if card1.Title.String() != cardTitle[0] {
		t.Fatal("Bad title [got]", card1.Title, "[want]", cardTitle[0])
	}
	if card1.Description.String() != desc[0] {
		t.Fatal("Bad description [got]", card1.Description, "[want]", desc[0])
	}

	// Update card
	card2, err := service.Update(gameID, collectionID, deckID, card1.ID, &UpdateCardDTO{
		Title:       cardTitle[1],
		Description: desc[1],
	})
	if err != nil {
		t.Fatal(err)
	}
	if card2.Title.String() != cardTitle[1] {
		t.Fatal("Bad title [got]", card2.Title, "[want]", cardTitle[1])
	}
	if card2.Description.String() != desc[1] {
		t.Fatal("Bad description [got]", card2.Description, "[want]", desc[1])
	}

	// Delete card
	err = service.Delete(gameID, collectionID, deckID, card2.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to update non-existing card
	_, err = service.Update(gameID, collectionID, deckID, card2.ID, &UpdateCardDTO{})
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}
}
func testList(t *testing.T) {
	service := NewService()
	cardTitle := []string{"B card", "A card"}

	// Empty list
	items, err := service.List(gameID, collectionID, deckID, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}

	// Create first card
	card1, err := service.Create(gameID, collectionID, deckID, &CreateCardDTO{
		Title: cardTitle[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// One card
	items, err = service.List(gameID, collectionID, deckID, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatal("List should be with 1 element")
	}

	// Create second card
	card2, err := service.Create(gameID, collectionID, deckID, &CreateCardDTO{
		Title: cardTitle[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sort by name
	items, err = service.List(gameID, collectionID, deckID, "name")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Title.String() != cardTitle[1] {
		t.Fatal("Bad name order: [got]", items[0].Title, "[want]", cardTitle[1])
	}
	if items[1].Title.String() != cardTitle[0] {
		t.Fatal("Bad name order: [got]", items[1].Title, "[want]", cardTitle[0])
	}

	// Sort by name_desc
	items, err = service.List(gameID, collectionID, deckID, "name_desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Title.String() != cardTitle[0] {
		t.Fatal("Bad name order: [got]", items[0].Title, "[want]", cardTitle[0])
	}
	if items[1].Title.String() != cardTitle[1] {
		t.Fatal("Bad name order: [got]", items[1].Title, "[want]", cardTitle[1])
	}

	// Sort by created date
	items, err = service.List(gameID, collectionID, deckID, "created")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Title.String() != cardTitle[0] {
		t.Fatal("Bad name order: [got]", items[0].Title, "[want]", cardTitle[0])
	}
	if items[1].Title.String() != cardTitle[1] {
		t.Fatal("Bad name order: [got]", items[1].Title, "[want]", cardTitle[1])
	}

	// Sort by created_desc
	items, err = service.List(gameID, collectionID, deckID, "created_desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Title.String() != cardTitle[1] {
		t.Fatal("Bad name order: [got]", items[0].Title, "[want]", cardTitle[1])
	}
	if items[1].Title.String() != cardTitle[0] {
		t.Fatal("Bad name order: [got]", items[1].Title, "[want]", cardTitle[0])
	}

	// Delete first card
	err = service.Delete(gameID, collectionID, deckID, card1.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Delete second card
	err = service.Delete(gameID, collectionID, deckID, card2.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Empty list
	items, err = service.List(gameID, collectionID, deckID, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}
}
func testItem(t *testing.T) {
	service := NewService()
	cardTitle := []string{"one", "two"}

	// Try to get non-existing card
	_, err := service.Item(gameID, collectionID, deckID, 1)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card1, err := service.Create(gameID, collectionID, deckID, &CreateCardDTO{
		Title: cardTitle[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid card
	_, err = service.Item(gameID, collectionID, deckID, card1.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid card
	_, err = service.Item(gameID, collectionID, deckID, 2)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Rename card
	card2, err := service.Update(gameID, collectionID, deckID, card1.ID, &UpdateCardDTO{Title: cardTitle[1]})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid card
	_, err = service.Item(gameID, collectionID, deckID, card1.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid card
	_, err = service.Item(gameID, collectionID, deckID, card2.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Delete card
	err = service.Delete(gameID, collectionID, deckID, card2.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCard(t *testing.T) {
	t.Parallel()

	// Set path for the game test artifacts
	dataPath := os.Getenv("TEST_DATA_PATH")
	if dataPath == "" {
		t.Fatal("TEST_DATA_PATH must be set")
	}
	config.GetConfig().SetDataPath(filepath.Join(dataPath, "card_test"))

	service := NewService()

	// Game not exist error
	_, err := service.Create(gameID, collectionID, deckID, &CreateCardDTO{
		Title: "test",
	})
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	gameService := games.NewService()
	_, err = gameService.Create(&games.CreateGameDTO{
		Name: gameID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Collection not exist error
	_, err = service.Create(gameID, collectionID, deckID, &CreateCardDTO{
		Title: "test",
	})
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}

	// Create collection
	collectionService := collections.NewService()
	_, err = collectionService.Create(gameID, &collections.CreateCollectionDTO{
		Name: collectionID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Deck not exist error
	_, err = service.Create(gameID, collectionID, deckID, &CreateCardDTO{
		Title: "test",
	})
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	deckService := decks.NewService()
	_, err = deckService.Create(gameID, collectionID, &decks.CreateDeckDTO{
		Type: deckID,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("create", testCreate)
	t.Run("delete", testDelete)
	t.Run("update", testUpdate)
	t.Run("list", testList)
	t.Run("item", testItem)
}

func fuzzCleanup(path string) {
	_ = os.RemoveAll(path)
}
func fuzzList(t *testing.T, service *CardService, waitItems int) error {
	items, err := service.List(gameID, collectionID, deckID, "")
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log("Get items error:", string(data))
		}
		return err
	}
	{
		data, _ := json.MarshalIndent(items, "", "	")
		t.Log("items:", string(data))
	}
	if len(items) != waitItems {
		{
			data, _ := json.MarshalIndent(items, "", "	")
			t.Log(string(data))
		}
		return fmt.Errorf("items: [wait] %d, [got] %d", waitItems, len(items))
	}
	return nil
}
func fuzzItem(t *testing.T, service *CardService, cardID int64, name, desc string) error {
	collection, err := service.Item(gameID, collectionID, deckID, cardID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	if collection.Title.String() != name {
		{
			data, _ := json.MarshalIndent(collection, "", "	")
			t.Log(string(data))
		}
		return fmt.Errorf("title: [wait] %s [got] %s", name, collection.Title)
	}
	if collection.Description.String() != desc {
		{
			data, _ := json.MarshalIndent(collection, "", "	")
			t.Log("item:", string(data))
		}
		return fmt.Errorf("description: [wait] %q [got] %q", desc, collection.Description)
	}
	return nil
}
func fuzzCreate(t *testing.T, service *CardService, name, desc string) (*CardInfo, error) {
	card, err := service.Create(gameID, collectionID, deckID, &CreateCardDTO{
		Title:       name,
		Description: desc,
	})
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return nil, err
	}
	{
		data, _ := json.MarshalIndent(card, "", "	")
		t.Log("create:", string(data))
	}
	return card, nil
}
func fuzzUpdate(t *testing.T, service *CardService, cardID int64, name, desc string) (*CardInfo, error) {
	card, err := service.Update(gameID, collectionID, deckID, cardID, &UpdateCardDTO{
		Title:       name,
		Description: desc,
	})
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return nil, err
	}
	{
		data, _ := json.MarshalIndent(card, "", "	")
		t.Log("update:", string(data))
	}
	return card, nil
}
func fuzzDelete(t *testing.T, service *CardService, cardID int64) error {
	err := service.Delete(gameID, collectionID, deckID, cardID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	return nil
}

func FuzzCard(f *testing.F) {
	// Set path for the collection test artifacts
	dataPath := os.Getenv("TEST_DATA_PATH")
	if dataPath == "" {
		f.Fatal("TEST_DATA_PATH must be set")
	}
	config.GetConfig().SetDataPath(filepath.Join(dataPath, "card_fuzz"))

	gameService := games.NewService()
	collectionService := collections.NewService()
	deckService := decks.NewService()
	service := NewService()

	msync := sync.Mutex{}
	f.Fuzz(func(t *testing.T, name1, desc1, name2, desc2 string) {
		items, err := gameService.List("")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) == 0 {
			// Create game
			_, err := gameService.Create(&games.CreateGameDTO{
				Name: gameID,
			})
			if err != nil {
				f.Fatal(err)
			}

			// Create collection
			_, err = collectionService.Create(gameID, &collections.CreateCollectionDTO{
				Name: collectionID,
			})
			if err != nil {
				f.Fatal(err)
			}

			// Create deck
			_, err = deckService.Create(gameID, collectionID, &decks.CreateDeckDTO{
				Type: deckID,
			})
			if err != nil {
				f.Fatal(err)
			}
		}

		// Only one test at once
		msync.Lock()
		defer msync.Unlock()

		// Empty list
		err = fuzzList(t, service, 0)
		if err != nil {
			t.Fatal(err)
		}

		// Create card
		card1, err := fuzzCreate(t, service, name1, desc1)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// List with card
		err = fuzzList(t, service, 1)
		if err != nil {
			{
				data, _ := json.MarshalIndent(err, "", "	")
				t.Log(string(data))
			}
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = fuzzItem(t, service, card1.ID, name1, desc1)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Update collection
		card2, err := fuzzUpdate(t, service, card1.ID, name2, desc2)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// List with card
		err = fuzzList(t, service, 1)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = fuzzItem(t, service, card2.ID, name2, desc2)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Delete card
		err = fuzzDelete(t, service, card2.ID)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Empty list
		err = fuzzList(t, service, 0)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}
	})
}
