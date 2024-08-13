package card

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCard "github.com/HardDie/DeckBuilder/internal/db/card"
	dbCollection "github.com/HardDie/DeckBuilder/internal/db/collection"
	dbCore "github.com/HardDie/DeckBuilder/internal/db/core"
	dbDeck "github.com/HardDie/DeckBuilder/internal/db/deck"
	dbGame "github.com/HardDie/DeckBuilder/internal/db/game"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	repositoriesCard "github.com/HardDie/DeckBuilder/internal/repositories/card"
	repositoriesCollection "github.com/HardDie/DeckBuilder/internal/repositories/collection"
	repositoriesDeck "github.com/HardDie/DeckBuilder/internal/repositories/deck"
	repositoriesGame "github.com/HardDie/DeckBuilder/internal/repositories/game"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
	servicesDeck "github.com/HardDie/DeckBuilder/internal/services/deck"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
)

type cardTest struct {
	gameID, collectionID, deckID string
	cfg                          *config.Config
	core                         dbCore.Core

	serviceGame       servicesGame.Game
	serviceCollection servicesCollection.Collection
	serviceDeck       servicesDeck.Deck
	serviceCard       Card
}

func newCardTest(t testing.TB) *cardTest {
	dir, err := os.MkdirTemp("", "card_test")
	if err != nil {
		t.Fatal("error creating temp dir", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	cfg := config.Get(false, "")
	cfg.SetDataPath(dir)

	fs := fsentry.NewFSEntry(cfg.Games())

	core := dbCore.New(fs)
	game := dbGame.New(fs)
	collection := dbCollection.New(fs, game)
	deck := dbDeck.New(fs, collection)
	card := dbCard.New(fs, deck)

	repositoryGame := repositoriesGame.New(cfg, game)
	repositoryCollection := repositoriesCollection.New(cfg, collection)
	repositoryDeck := repositoriesDeck.New(cfg, collection, deck)
	repositoryCard := repositoriesCard.New(cfg, card)

	return &cardTest{
		gameID:       "test_card__game",
		collectionID: "test_card__collection",
		deckID:       "test_card__deck",
		cfg:          cfg,
		core:         core,

		serviceGame:       servicesGame.New(cfg, repositoryGame),
		serviceCollection: servicesCollection.New(cfg, repositoryCollection),
		serviceDeck:       servicesDeck.New(cfg, repositoryDeck),
		serviceCard:       New(cfg, repositoryCard),
	}
}

func (tt *cardTest) testCreate(t *testing.T) {
	cardName := "create_one"
	desc := "best card ever"
	count := 2

	// Create card
	card, err := tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID+"_create", &dto.CreateCardDTO{
		Name:        cardName,
		Description: desc,
		Count:       count,
	})
	if err != nil {
		t.Fatal(err)
	}
	if card.Name != cardName {
		t.Fatal("Bad name [got]", card.Name, "[want]", cardName)
	}
	if card.Description != desc {
		t.Fatal("Bad description [got]", card.Description, "[want]", desc)
	}
	if card.Count != count {
		t.Fatal("Bad count [got]", card.Count, "[want]", count)
	}

	// Delete card
	err = tt.serviceCard.Delete(tt.gameID, tt.collectionID, tt.deckID+"_create", card.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *cardTest) testDelete(t *testing.T) {
	cardName := "delete_one"

	// Try to remove non-existing card
	err := tt.serviceCard.Delete(tt.gameID, tt.collectionID, tt.deckID+"_delete", 1)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card, err := tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID+"_delete", &dto.CreateCardDTO{
		Name: cardName,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete card
	err = tt.serviceCard.Delete(tt.gameID, tt.collectionID, tt.deckID+"_delete", card.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete card twice
	err = tt.serviceCard.Delete(tt.gameID, tt.collectionID, tt.deckID+"_delete", card.ID)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}
}
func (tt *cardTest) testUpdate(t *testing.T) {
	cardName := []string{"update_one", "update_two"}
	desc := []string{"first description", "second description"}
	count := []int{5, 12}

	// Try to update non-existing card
	_, err := tt.serviceCard.Update(tt.gameID, tt.collectionID, tt.deckID+"_update", 1, &dto.UpdateCardDTO{})
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card1, err := tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID+"_update", &dto.CreateCardDTO{
		Name:        cardName[0],
		Description: desc[0],
		Count:       count[0],
	})
	if err != nil {
		t.Fatal(err)
	}
	if card1.Name != cardName[0] {
		t.Fatal("Bad name [got]", card1.Name, "[want]", cardName[0])
	}
	if card1.Description != desc[0] {
		t.Fatal("Bad description [got]", card1.Description, "[want]", desc[0])
	}
	if card1.Count != count[0] {
		t.Fatal("Bad count [got]", card1.Count, "[want]", count[0])
	}

	// Update card
	card2, err := tt.serviceCard.Update(tt.gameID, tt.collectionID, tt.deckID+"_update", card1.ID, &dto.UpdateCardDTO{
		Name:        cardName[1],
		Description: desc[1],
		Count:       count[1],
	})
	if err != nil {
		t.Fatal(err)
	}
	if card2.Name != cardName[1] {
		t.Fatal("Bad name [got]", card2.Name, "[want]", cardName[1])
	}
	if card2.Description != desc[1] {
		t.Fatal("Bad description [got]", card2.Description, "[want]", desc[1])
	}
	if card2.Count != count[1] {
		t.Fatal("Bad count [got]", card2.Count, "[want]", count[1])
	}

	// Delete card
	err = tt.serviceCard.Delete(tt.gameID, tt.collectionID, tt.deckID+"_update", card2.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to update non-existing card
	_, err = tt.serviceCard.Update(tt.gameID, tt.collectionID, tt.deckID+"_update", card2.ID, &dto.UpdateCardDTO{})
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
	items, _, err := tt.serviceCard.List(tt.gameID, tt.collectionID, tt.deckID+"_list", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}

	// Create first card
	card1, err := tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID+"_list", &dto.CreateCardDTO{
		Name: cardName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// One card
	items, _, err = tt.serviceCard.List(tt.gameID, tt.collectionID, tt.deckID+"_list", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatal("List should be with 1 element")
	}

	// Create second card
	card2, err := tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID+"_list", &dto.CreateCardDTO{
		Name: cardName[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sort by name
	items, _, err = tt.serviceCard.List(tt.gameID, tt.collectionID, tt.deckID+"_list", "name", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != cardName[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", cardName[1])
	}
	if items[1].Name != cardName[0] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", cardName[0])
	}

	// Sort by name_desc
	items, _, err = tt.serviceCard.List(tt.gameID, tt.collectionID, tt.deckID+"_list", "name_desc", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != cardName[0] {
		t.Fatal("Bad name_desc order: [got]", items[0].Name, "[want]", cardName[0])
	}
	if items[1].Name != cardName[1] {
		t.Fatal("Bad name_desc order: [got]", items[1].Name, "[want]", cardName[1])
	}

	// Sort by created date
	items, _, err = tt.serviceCard.List(tt.gameID, tt.collectionID, tt.deckID+"_list", "created", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != cardName[0] {
		t.Fatal("Bad created order: [got]", items[0].Name, "[want]", cardName[0])
	}
	if items[1].Name != cardName[1] {
		t.Fatal("Bad created order: [got]", items[1].Name, "[want]", cardName[1])
	}

	// Sort by created_desc
	items, _, err = tt.serviceCard.List(tt.gameID, tt.collectionID, tt.deckID+"_list", "created_desc", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != cardName[1] {
		t.Fatal("Bad created_desc order: [got]", items[0].Name, "[want]", cardName[1])
	}
	if items[1].Name != cardName[0] {
		t.Fatal("Bad created_desc order: [got]", items[1].Name, "[want]", cardName[0])
	}

	// Delete first card
	err = tt.serviceCard.Delete(tt.gameID, tt.collectionID, tt.deckID+"_list", card1.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Delete second card
	err = tt.serviceCard.Delete(tt.gameID, tt.collectionID, tt.deckID+"_list", card2.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Empty list
	items, _, err = tt.serviceCard.List(tt.gameID, tt.collectionID, tt.deckID+"_list", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}
}
func (tt *cardTest) testItem(t *testing.T) {
	cardName := []string{"item_one", "item_two"}

	// Try to get non-existing card
	_, err := tt.serviceCard.Item(tt.gameID, tt.collectionID, tt.deckID+"_item", 1)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card1, err := tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID+"_item", &dto.CreateCardDTO{
		Name: cardName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid card
	_, err = tt.serviceCard.Item(tt.gameID, tt.collectionID, tt.deckID+"_item", card1.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid card
	_, err = tt.serviceCard.Item(tt.gameID, tt.collectionID, tt.deckID+"_item", 2)
	if err == nil {
		t.Fatal("Error, card not exist")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Rename card
	card2, err := tt.serviceCard.Update(tt.gameID, tt.collectionID, tt.deckID+"_item", card1.ID, &dto.UpdateCardDTO{Name: cardName[1]})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid card
	_, err = tt.serviceCard.Item(tt.gameID, tt.collectionID, tt.deckID+"_item", card1.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid card
	_, err = tt.serviceCard.Item(tt.gameID, tt.collectionID, tt.deckID+"_item", card2.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Delete card
	err = tt.serviceCard.Delete(tt.gameID, tt.collectionID, tt.deckID+"_item", card2.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *cardTest) testImage(t *testing.T) {
	cardTitle := "image_one"
	pngImage := "https://github.com/fluidicon.png"
	jpegImage := "https://avatars.githubusercontent.com/apple"

	// Check no card
	_, _, err := tt.serviceCard.GetImage(tt.gameID, tt.collectionID, tt.deckID+"_image", 1)
	if err == nil {
		t.Fatal("Error, card not exists")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card, err := tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID+"_image", &dto.CreateCardDTO{
		Name:  cardTitle,
		Image: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.serviceCard.GetImage(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update card
	_, err = tt.serviceCard.Update(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID, &dto.UpdateCardDTO{
		Name:  cardTitle,
		Image: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceCard.GetImage(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update card
	_, err = tt.serviceCard.Update(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID, &dto.UpdateCardDTO{
		Name:  cardTitle,
		Image: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.serviceCard.GetImage(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID)
	if err == nil {
		t.Fatal("Error, card don't have image")
	}
	if !errors.Is(err, er.CardImageNotExists) {
		t.Fatal(err)
	}

	// Delete card
	err = tt.serviceCard.Delete(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *cardTest) testImageBin(t *testing.T) {
	cardTitle := "image_bin_one"

	pageImage := images.CreateImage(100, 100)
	pngImage, err := images.ImageToPng(pageImage)
	if err != nil {
		t.Fatal(err)
	}
	jpegImage, err := images.ImageToJpeg(pageImage)
	if err != nil {
		t.Fatal(err)
	}
	gifImage, err := images.ImageToGif(pageImage)
	if err != nil {
		t.Fatal(err)
	}

	// Check no card
	_, _, err = tt.serviceCard.GetImage(tt.gameID, tt.collectionID, tt.deckID+"_image", 1)
	if err == nil {
		t.Fatal("Error, card not exists")
	}
	if !errors.Is(err, er.CardNotExists) {
		t.Fatal(err)
	}

	// Create card
	card, err := tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID+"_image", &dto.CreateCardDTO{
		Name:      cardTitle,
		ImageFile: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.serviceCard.GetImage(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update card
	_, err = tt.serviceCard.Update(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID, &dto.UpdateCardDTO{
		Name:      cardTitle,
		ImageFile: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceCard.GetImage(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update card
	_, err = tt.serviceCard.Update(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID, &dto.UpdateCardDTO{
		Name:      cardTitle,
		ImageFile: gifImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceCard.GetImage(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "gif" {
		t.Fatal("Image type error! [got]", imgType, "[want] gif")
	}

	// Update card
	_, err = tt.serviceCard.Update(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID, &dto.UpdateCardDTO{
		Name: cardTitle,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceCard.GetImage(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "gif" {
		t.Fatal("Image type error! [got]", imgType, "[want] gif")
	}

	// Update card
	_, err = tt.serviceCard.Update(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID, &dto.UpdateCardDTO{
		Name:  cardTitle,
		Image: "empty",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.serviceCard.GetImage(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID)
	if err == nil {
		t.Fatal("Error, card don't have image")
	}
	if !errors.Is(err, er.CardImageNotExists) {
		t.Fatal(err)
	}

	// Delete card
	err = tt.serviceCard.Delete(tt.gameID, tt.collectionID, tt.deckID+"_image", card.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCard(t *testing.T) {
	t.Parallel()

	tt := newCardTest(t)

	if err := tt.core.Init(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := tt.core.Drop(); err != nil {
			t.Fatal(err)
		}
	}()

	// Game not exist error
	_, err := tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name: "test",
	})
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = tt.serviceGame.Create(&dto.CreateGameDTO{
		Name: tt.gameID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Collection not exist error
	_, err = tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name: "test",
	})
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}

	// Create collection
	_, err = tt.serviceCollection.Create(tt.gameID, &dto.CreateCollectionDTO{
		Name: tt.collectionID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Deck not exist error
	_, err = tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name: "test",
	})
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	decks := []string{"_create", "_delete", "_update", "_list", "_item", "_image"}
	for _, deck := range decks {
		// Create deck
		_, err = tt.serviceDeck.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
			Name: tt.deckID + deck,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("create", tt.testCreate)
	t.Run("delete", tt.testDelete)
	t.Run("update", tt.testUpdate)
	t.Run("list", tt.testList)
	t.Run("item", tt.testItem)
	t.Run("image", tt.testImage)
	t.Run("image_bin", tt.testImageBin)
}

func (tt *cardTest) fuzzCleanup() {
	_ = tt.core.Drop()
	_ = tt.core.Init()
}
func (tt *cardTest) fuzzList(t *testing.T, waitItems int) error {
	items, _, err := tt.serviceCard.List(tt.gameID, tt.collectionID, tt.deckID, "", "")
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
	card, err := tt.serviceCard.Item(tt.gameID, tt.collectionID, tt.deckID, cardID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	if card.Name != name {
		{
			data, _ := json.MarshalIndent(card, "", "	")
			t.Log(string(data))
		}
		return fmt.Errorf("title: [wait] %s [got] %s", name, card.Name)
	}
	if card.Description != desc {
		{
			data, _ := json.MarshalIndent(card, "", "	")
			t.Log("item:", string(data))
		}
		return fmt.Errorf("description: [wait] %q [got] %q", desc, card.Description)
	}
	if val, ok := card.Variables[name]; !ok {
		{
			data, _ := json.MarshalIndent(card, "", "	")
			t.Log("item:", string(data))
		}
		return fmt.Errorf("value key not found: %q", name)
	} else if val != desc {
		{
			data, _ := json.MarshalIndent(card, "", "	")
			t.Log("item:", string(data))
		}
		return fmt.Errorf("value: [wait] %q [got] %q", desc, val)
	}
	return nil
}
func (tt *cardTest) fuzzCreate(t *testing.T, name, desc string) (*entity.CardInfo, error) {
	card, err := tt.serviceCard.Create(tt.gameID, tt.collectionID, tt.deckID, &dto.CreateCardDTO{
		Name:        name,
		Description: desc,
		Variables: map[string]string{
			name: desc,
		},
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
	card, err := tt.serviceCard.Update(tt.gameID, tt.collectionID, tt.deckID, cardID, &dto.UpdateCardDTO{
		Name:        name,
		Description: desc,
		Variables: map[string]string{
			name: desc,
		},
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
	err := tt.serviceCard.Delete(tt.gameID, tt.collectionID, tt.deckID, cardID)
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
	tt := newCardTest(f)

	if err := tt.core.Init(); err != nil {
		f.Fatal(err)
	}
	defer func() {
		if err := tt.core.Drop(); err != nil {
			f.Fatal(err)
		}
	}()

	f.Fuzz(func(t *testing.T, name1, desc1, name2, desc2 string) {
		items, _, err := tt.serviceGame.List("", "")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) == 0 {
			// Create game
			_, err := tt.serviceGame.Create(&dto.CreateGameDTO{
				Name: tt.gameID,
			})
			if err != nil {
				f.Fatal(err)
			}

			// Create collection
			_, err = tt.serviceCollection.Create(tt.gameID, &dto.CreateCollectionDTO{
				Name: tt.collectionID,
			})
			if err != nil {
				f.Fatal(err)
			}

			// Create deck
			_, err = tt.serviceDeck.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
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
