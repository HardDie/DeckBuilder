package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/repository"
)

type cardTest struct {
	gameID, collectionID, deckID string
	cfg                          *config.Config
	gameService                  IGameService
	collectionService            ICollectionService
	deckService                  IDeckService
	cardService                  ICardService
}

func newCardTest(dataPath string) *cardTest {
	cfg := config.Get()
	cfg.SetDataPath(dataPath)

	gameRepository := repository.NewGameRepository(cfg)
	collectionRepository := repository.NewCollectionRepository(cfg, gameRepository)
	deckRepository := repository.NewDeckRepository(cfg, collectionRepository)

	return &cardTest{
		gameID:            "test_card__game",
		collectionID:      "test_card__collection",
		deckID:            "test_card__deck",
		cfg:               cfg,
		gameService:       NewGameService(gameRepository),
		collectionService: NewCollectionService(collectionRepository),
		deckService:       NewDeckService(deckRepository),
		cardService:       NewCardService(repository.NewCardRepository(cfg, deckRepository)),
	}
}

func (tt *cardTest) testCreate(t *testing.T) {
	cardName := "one"
	desc := "best card ever"
	count := 2

	// Create card
	card, err := tt.cardService.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name:        cardName,
		Description: desc,
		Count:       count,
	})
	if err != nil {
		t.Fatal(err)
	}
	if card.Name.String() != cardName {
		t.Fatal("Bad name [got]", card.Name, "[want]", cardName)
	}
	if card.Description.String() != desc {
		t.Fatal("Bad description [got]", card.Description, "[want]", desc)
	}
	if card.Count != count {
		t.Fatal("Bad count [got]", card.Count, "[want]", count)
	}

	// Delete card
	err = tt.cardService.Delete(tt.gameID, tt.collectionID, tt.deckID, card.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *cardTest) testDelete(t *testing.T) {
	cardName := "one"

	// Try to remove non-existing card
	err := tt.cardService.Delete(tt.gameID, tt.collectionID, tt.deckID, 1)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card, err := tt.cardService.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name: cardName,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete card
	err = tt.cardService.Delete(tt.gameID, tt.collectionID, tt.deckID, card.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete card twice
	err = tt.cardService.Delete(tt.gameID, tt.collectionID, tt.deckID, card.ID)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}
}
func (tt *cardTest) testUpdate(t *testing.T) {
	cardName := []string{"one", "two"}
	desc := []string{"first description", "second description"}
	count := []int{5, 12}

	// Try to update non-existing card
	_, err := tt.cardService.Update(tt.gameID, tt.collectionID, tt.deckID, 1, &dto.UpdateCardDTO{})
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card1, err := tt.cardService.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name:        cardName[0],
		Description: desc[0],
		Count:       count[0],
	})
	if err != nil {
		t.Fatal(err)
	}
	if card1.Name.String() != cardName[0] {
		t.Fatal("Bad name [got]", card1.Name, "[want]", cardName[0])
	}
	if card1.Description.String() != desc[0] {
		t.Fatal("Bad description [got]", card1.Description, "[want]", desc[0])
	}
	if card1.Count != count[0] {
		t.Fatal("Bad count [got]", card1.Count, "[want]", count[0])
	}

	// Update card
	card2, err := tt.cardService.Update(tt.gameID, tt.collectionID, tt.deckID, card1.ID, &dto.UpdateCardDTO{
		Name:        cardName[1],
		Description: desc[1],
		Count:       count[1],
	})
	if err != nil {
		t.Fatal(err)
	}
	if card2.Name.String() != cardName[1] {
		t.Fatal("Bad name [got]", card2.Name, "[want]", cardName[1])
	}
	if card2.Description.String() != desc[1] {
		t.Fatal("Bad description [got]", card2.Description, "[want]", desc[1])
	}
	if card2.Count != count[1] {
		t.Fatal("Bad count [got]", card2.Count, "[want]", count[1])
	}

	// Delete card
	err = tt.cardService.Delete(tt.gameID, tt.collectionID, tt.deckID, card2.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to update non-existing card
	_, err = tt.cardService.Update(tt.gameID, tt.collectionID, tt.deckID, card2.ID, &dto.UpdateCardDTO{})
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}
}
func (tt *cardTest) testList(t *testing.T) {
	cardName := []string{"B card", "A card"}

	// Empty list
	items, err := tt.cardService.List(tt.gameID, tt.collectionID, tt.deckID, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}

	// Create first card
	card1, err := tt.cardService.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name: cardName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// One card
	items, err = tt.cardService.List(tt.gameID, tt.collectionID, tt.deckID, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatal("List should be with 1 element")
	}

	// Create second card
	card2, err := tt.cardService.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name: cardName[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sort by name
	items, err = tt.cardService.List(tt.gameID, tt.collectionID, tt.deckID, "name")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != cardName[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", cardName[1])
	}
	if items[1].Name.String() != cardName[0] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", cardName[0])
	}

	// Sort by name_desc
	items, err = tt.cardService.List(tt.gameID, tt.collectionID, tt.deckID, "name_desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != cardName[0] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", cardName[0])
	}
	if items[1].Name.String() != cardName[1] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", cardName[1])
	}

	// Sort by created date
	items, err = tt.cardService.List(tt.gameID, tt.collectionID, tt.deckID, "created")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != cardName[0] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", cardName[0])
	}
	if items[1].Name.String() != cardName[1] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", cardName[1])
	}

	// Sort by created_desc
	items, err = tt.cardService.List(tt.gameID, tt.collectionID, tt.deckID, "created_desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != cardName[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", cardName[1])
	}
	if items[1].Name.String() != cardName[0] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", cardName[0])
	}

	// Delete first card
	err = tt.cardService.Delete(tt.gameID, tt.collectionID, tt.deckID, card1.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Delete second card
	err = tt.cardService.Delete(tt.gameID, tt.collectionID, tt.deckID, card2.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Empty list
	items, err = tt.cardService.List(tt.gameID, tt.collectionID, tt.deckID, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}
}
func (tt *cardTest) testItem(t *testing.T) {
	cardName := []string{"one", "two"}

	// Try to get non-existing card
	_, err := tt.cardService.Item(tt.gameID, tt.collectionID, tt.deckID, 1)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card1, err := tt.cardService.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name: cardName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid card
	_, err = tt.cardService.Item(tt.gameID, tt.collectionID, tt.deckID, card1.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid card
	_, err = tt.cardService.Item(tt.gameID, tt.collectionID, tt.deckID, 2)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Rename card
	card2, err := tt.cardService.Update(tt.gameID, tt.collectionID, tt.deckID, card1.ID, &dto.UpdateCardDTO{Name: cardName[1]})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid card
	_, err = tt.cardService.Item(tt.gameID, tt.collectionID, tt.deckID, card1.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid card
	_, err = tt.cardService.Item(tt.gameID, tt.collectionID, tt.deckID, card2.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Delete card
	err = tt.cardService.Delete(tt.gameID, tt.collectionID, tt.deckID, card2.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *cardTest) testImage(t *testing.T) {
	cardTitle := "one"
	pngImage := "https://github.com/fluidicon.png"
	jpegImage := "https://avatars.githubusercontent.com/apple"

	// Check no card
	_, _, err := tt.cardService.GetImage(tt.gameID, tt.collectionID, tt.deckID, 1)
	if err == nil {
		t.Fatal("Error, card not exists")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card, err := tt.cardService.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name:  cardTitle,
		Image: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.cardService.GetImage(tt.gameID, tt.collectionID, tt.deckID, card.ID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update card
	_, err = tt.cardService.Update(tt.gameID, tt.collectionID, tt.deckID, card.ID, &dto.UpdateCardDTO{
		Image: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.cardService.GetImage(tt.gameID, tt.collectionID, tt.deckID, card.ID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update card
	_, err = tt.cardService.Update(tt.gameID, tt.collectionID, tt.deckID, card.ID, &dto.UpdateCardDTO{
		Image: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.cardService.GetImage(tt.gameID, tt.collectionID, tt.deckID, card.ID)
	if err == nil {
		t.Fatal("Error, card don't have image")
	}
	if !errors.Is(err, er.CardImageNotExists) {
		t.Fatal(err)
	}

	// Delete card
	err = tt.cardService.Delete(tt.gameID, tt.collectionID, tt.deckID, card.ID)
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
	tt := newCardTest(filepath.Join(dataPath, "card_test"))

	// Game not exist error
	_, err := tt.cardService.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name: "test",
	})
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = tt.gameService.Create(&dto.CreateGameDTO{
		Name: tt.gameID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Collection not exist error
	_, err = tt.cardService.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name: "test",
	})
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}

	// Create collection
	_, err = tt.collectionService.Create(tt.gameID, &dto.CreateCollectionDTO{
		Name: tt.collectionID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Deck not exist error
	_, err = tt.cardService.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name: "test",
	})
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
		Name: tt.deckID,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("create", tt.testCreate)
	t.Run("delete", tt.testDelete)
	t.Run("update", tt.testUpdate)
	t.Run("list", tt.testList)
	t.Run("item", tt.testItem)
	t.Run("image", tt.testImage)
}

func (tt *cardTest) fuzzCleanup() {
	_ = os.RemoveAll(tt.cfg.Data)
}
func (tt *cardTest) fuzzList(t *testing.T, waitItems int) error {
	items, err := tt.cardService.List(tt.gameID, tt.collectionID, tt.deckID, "")
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
func (tt *cardTest) fuzzItem(t *testing.T, cardID int64, name, desc string) error {
	collection, err := tt.cardService.Item(tt.gameID, tt.collectionID, tt.deckID, cardID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	if collection.Name.String() != name {
		{
			data, _ := json.MarshalIndent(collection, "", "	")
			t.Log(string(data))
		}
		return fmt.Errorf("title: [wait] %s [got] %s", name, collection.Name)
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
func (tt *cardTest) fuzzCreate(t *testing.T, name, desc string) (*entity.CardInfo, error) {
	card, err := tt.cardService.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name:        name,
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
func (tt *cardTest) fuzzUpdate(t *testing.T, cardID int64, name, desc string) (*entity.CardInfo, error) {
	card, err := tt.cardService.Update(tt.gameID, tt.collectionID, tt.deckID, cardID, &dto.UpdateCardDTO{
		Name:        name,
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
func (tt *cardTest) fuzzDelete(t *testing.T, cardID int64) error {
	err := tt.cardService.Delete(tt.gameID, tt.collectionID, tt.deckID, cardID)
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
	tt := newCardTest(filepath.Join(dataPath, "card_fuzz_"+uuid.New().String()))

	f.Fuzz(func(t *testing.T, name1, desc1, name2, desc2 string) {
		items, err := tt.gameService.List("")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) == 0 {
			// Create game
			_, err := tt.gameService.Create(&dto.CreateGameDTO{
				Name: tt.gameID,
			})
			if err != nil {
				f.Fatal(err)
			}

			// Create collection
			_, err = tt.collectionService.Create(tt.gameID, &dto.CreateCollectionDTO{
				Name: tt.collectionID,
			})
			if err != nil {
				f.Fatal(err)
			}

			// Create deck
			_, err = tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
				Name: tt.deckID,
			})
			if err != nil {
				f.Fatal(err)
			}
		}

		// Empty list
		err = tt.fuzzList(t, 0)
		if err != nil {
			t.Fatal(err)
		}

		// Create card
		card1, err := tt.fuzzCreate(t, name1, desc1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// List with card
		err = tt.fuzzList(t, 1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = tt.fuzzItem(t, card1.ID, name1, desc1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Update collection
		card2, err := tt.fuzzUpdate(t, card1.ID, name2, desc2)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// List with card
		err = tt.fuzzList(t, 1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = tt.fuzzItem(t, card2.ID, name2, desc2)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Delete card
		err = tt.fuzzDelete(t, card2.ID)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Empty list
		err = tt.fuzzList(t, 0)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}
	})
}
