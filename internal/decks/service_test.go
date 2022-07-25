package decks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/google/uuid"

	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
	er "tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
	"tts_deck_build/internal/utils"
)

var (
	gameID       = "test_game"
	collectionID = "test_collection"
)

func testCreate(t *testing.T) {
	service := NewService()
	deckType := "one"

	// Create deck
	deck, err := service.Create(gameID, collectionID, &CreateDeckDTO{
		Type: deckType,
	})
	if err != nil {
		t.Fatal(err)
	}
	if deck.Type.String() != deckType {
		t.Fatal("Bad type [got]", deck.Type, "[want]", deckType)
	}

	// Try to create duplicate
	_, err = service.Create(gameID, collectionID, &CreateDeckDTO{
		Type: deckType,
	})
	if err == nil {
		t.Fatal("Error, you can't create duplicate deck")
	}
	if !errors.Is(err, er.DeckExist) {
		t.Fatal(err)
	}

	// Delete deck
	err = service.Delete(gameID, collectionID, deck.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func testDelete(t *testing.T) {
	service := NewService()
	deckType := "one"
	deckID := utils.NameToID(deckType)

	// Try to remove non-existing deck
	err := service.Delete(gameID, collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = service.Create(gameID, collectionID, &CreateDeckDTO{
		Type: deckType,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete deck
	err = service.Delete(gameID, collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete deck twice
	err = service.Delete(gameID, collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}
}
func testUpdate(t *testing.T) {
	service := NewService()
	deckType := []string{"one", "two"}
	deckID := []string{utils.NameToID(deckType[0]), utils.NameToID(deckType[1])}

	// Try to update non-existing deck
	_, err := service.Update(gameID, collectionID, deckID[0], &UpdateDeckDTO{})
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	deck, err := service.Create(gameID, collectionID, &CreateDeckDTO{
		Type: deckType[0],
	})
	if err != nil {
		t.Fatal(err)
	}
	if deck.Type.String() != deckType[0] {
		t.Fatal("Bad type [got]", deck.Type, "[want]", deckType[0])
	}

	// Update deck
	deck, err = service.Update(gameID, collectionID, deckID[0], &UpdateDeckDTO{
		Type: deckType[1],
	})
	if err != nil {
		t.Fatal(err)
	}
	if deck.Type.String() != deckType[1] {
		t.Fatal("Bad type [got]", deck.Type, "[want]", deckType[1])
	}

	// Delete deck
	err = service.Delete(gameID, collectionID, deckID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Try to update non-existing deck
	_, err = service.Update(gameID, collectionID, deckID[1], &UpdateDeckDTO{})
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}
}
func testList(t *testing.T) {
	service := NewService()
	deckType := []string{"B deck", "A deck"}
	deckID := []string{utils.NameToID(deckType[0]), utils.NameToID(deckType[1])}

	// Empty list
	items, err := service.List(gameID, collectionID, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}

	// Create first deck
	_, err = service.Create(gameID, collectionID, &CreateDeckDTO{
		Type: deckType[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// One deck
	items, err = service.List(gameID, collectionID, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatal("List should be with 1 element")
	}

	// Create second deck
	_, err = service.Create(gameID, collectionID, &CreateDeckDTO{
		Type: deckType[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sort by name
	items, err = service.List(gameID, collectionID, "name")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Type.String() != deckType[1] {
		t.Fatal("Bad name order: [got]", items[0].Type, "[want]", deckType[1])
	}
	if items[1].Type.String() != deckType[0] {
		t.Fatal("Bad name order: [got]", items[1].Type, "[want]", deckType[0])
	}

	// Sort by name_desc
	items, err = service.List(gameID, collectionID, "name_desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Type.String() != deckType[0] {
		t.Fatal("Bad name order: [got]", items[0].Type, "[want]", deckType[0])
	}
	if items[1].Type.String() != deckType[1] {
		t.Fatal("Bad name order: [got]", items[1].Type, "[want]", deckType[1])
	}

	// Sort by created date
	items, err = service.List(gameID, collectionID, "created")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Type.String() != deckType[0] {
		t.Fatal("Bad name order: [got]", items[0].Type, "[want]", deckType[0])
	}
	if items[1].Type.String() != deckType[1] {
		t.Fatal("Bad name order: [got]", items[1].Type, "[want]", deckType[1])
	}

	// Sort by created_desc
	items, err = service.List(gameID, collectionID, "created_desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Type.String() != deckType[1] {
		t.Fatal("Bad name order: [got]", items[0].Type, "[want]", deckType[1])
	}
	if items[1].Type.String() != deckType[0] {
		t.Fatal("Bad name order: [got]", items[1].Type, "[want]", deckType[0])
	}

	// Delete first deck
	err = service.Delete(gameID, collectionID, deckID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Delete second deck
	err = service.Delete(gameID, collectionID, deckID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Empty list
	items, err = service.List(gameID, collectionID, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}
}
func testItem(t *testing.T) {
	service := NewService()
	deckType := []string{"one", "two"}
	deckID := []string{utils.NameToID(deckType[0]), utils.NameToID(deckType[1])}

	// Try to get non-existing deck
	_, err := service.Item(gameID, collectionID, deckID[0])
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = service.Create(gameID, collectionID, &CreateDeckDTO{
		Type: deckType[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid deck
	_, err = service.Item(gameID, collectionID, deckID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid deck
	_, err = service.Item(gameID, collectionID, deckID[1])
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Rename deck
	_, err = service.Update(gameID, collectionID, deckID[0], &UpdateDeckDTO{Type: deckType[1]})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid deck
	_, err = service.Item(gameID, collectionID, deckID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid deck
	_, err = service.Item(gameID, collectionID, deckID[0])
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Delete deck
	err = service.Delete(gameID, collectionID, deckID[1])
	if err != nil {
		t.Fatal(err)
	}
}
func testImage(t *testing.T) {
	service := NewService()
	deckType := "one"
	deckID := utils.NameToID(deckType)
	pngImage := "https://github.com/fluidicon.png"
	jpegImage := "https://avatars.githubusercontent.com/apple"

	// Check no deck
	_, _, err := service.GetImage(gameID, collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck not exists")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = service.Create(gameID, collectionID, &CreateDeckDTO{
		Type:          deckType,
		BacksideImage: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := service.GetImage(gameID, collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update deck
	_, err = service.Update(gameID, collectionID, deckID, &UpdateDeckDTO{
		BacksideImage: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = service.GetImage(gameID, collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update deck
	_, err = service.Update(gameID, collectionID, deckID, &UpdateDeckDTO{
		BacksideImage: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = service.GetImage(gameID, collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck don't have image")
	}
	if !errors.Is(err, er.DeckImageNotExists) {
		t.Fatal(err)
	}

	// Delete deck
	err = service.Delete(gameID, collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeck(t *testing.T) {
	t.Parallel()

	// Set path for the game test artifacts
	dataPath := os.Getenv("TEST_DATA_PATH")
	if dataPath == "" {
		t.Fatal("TEST_DATA_PATH must be set")
	}
	config.GetConfig().SetDataPath(filepath.Join(dataPath, "deck_test"))

	service := NewService()

	// Game not exist error
	_, err := service.Create(gameID, collectionID, &CreateDeckDTO{
		Type: "test",
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
	_, err = service.Create(gameID, collectionID, &CreateDeckDTO{
		Type: "test",
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

	t.Run("create", testCreate)
	t.Run("delete", testDelete)
	t.Run("update", testUpdate)
	t.Run("list", testList)
	t.Run("item", testItem)
	t.Run("image", testImage)
}

func fuzzCleanup(path string) {
	_ = os.RemoveAll(path)
}
func fuzzList(t *testing.T, service *DeckService, waitItems int) error {
	items, err := service.List(gameID, collectionID, "")
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
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
func fuzzItem(t *testing.T, service *DeckService, deckID, deckType string) error {
	deck, err := service.Item(gameID, collectionID, deckID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	if deck.Type.String() != deckType {
		{
			data, _ := json.MarshalIndent(deck, "", "	")
			t.Log(string(data))
		}
		return fmt.Errorf("name: [wait] %s [got] %s", deckType, deck.Type)
	}
	return nil
}
func fuzzCreate(t *testing.T, service *DeckService, deckType string) (*DeckInfo, error) {
	deck, err := service.Create(gameID, collectionID, &CreateDeckDTO{
		Type: deckType,
	})
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return nil, err
	}
	{
		data, _ := json.MarshalIndent(deck, "", "	")
		t.Log("create:", string(data))
	}
	return deck, nil
}
func fuzzUpdate(t *testing.T, service *DeckService, deckID, deckType string) (*DeckInfo, error) {
	deck, err := service.Update(gameID, collectionID, deckID, &UpdateDeckDTO{
		Type: deckType,
	})
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return nil, err
	}
	{
		data, _ := json.MarshalIndent(deck, "", "	")
		t.Log("update:", string(data))
	}
	return deck, nil
}
func fuzzDelete(t *testing.T, service *DeckService, deckID string) error {
	err := service.Delete(gameID, collectionID, deckID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	return nil
}

func FuzzDeck(f *testing.F) {
	// Set path for the deck test artifacts
	dataPath := os.Getenv("TEST_DATA_PATH")
	if dataPath == "" {
		f.Fatal("TEST_DATA_PATH must be set")
	}
	config.GetConfig().SetDataPath(filepath.Join(dataPath, "deck_fuzz_"+uuid.New().String()))

	gameService := games.NewService()
	collectionService := collections.NewService()
	service := NewService()

	msync := sync.Mutex{}
	f.Fuzz(func(t *testing.T, type1, type2 string) {
		gameItems, err := gameService.List("")
		if err != nil {
			t.Fatal(err)
		}
		if len(gameItems) == 0 {
			// Create game
			_, err = gameService.Create(&games.CreateGameDTO{
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
		}

		if utils.NameToID(type1) == "" || utils.NameToID(type2) == "" {
			// skip
			return
		}

		// Only one test at once
		msync.Lock()
		defer msync.Unlock()

		// Empty list
		err = fuzzList(t, service, 0)
		if err != nil {
			t.Fatal(err)
		}

		// Create deck
		deck1, err := fuzzCreate(t, service, type1)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// List with deck
		err = fuzzList(t, service, 1)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = fuzzItem(t, service, deck1.ID, type1)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Update collection
		collection2, err := fuzzUpdate(t, service, utils.NameToID(type1), type2)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// List with collection
		err = fuzzList(t, service, 1)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = fuzzItem(t, service, collection2.ID, type2)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Delete collection
		err = fuzzDelete(t, service, utils.NameToID(type2))
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
