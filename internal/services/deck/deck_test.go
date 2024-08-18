package deck

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCollection "github.com/HardDie/DeckBuilder/internal/db/collection"
	dbCore "github.com/HardDie/DeckBuilder/internal/db/core"
	dbDeck "github.com/HardDie/DeckBuilder/internal/db/deck"
	dbGame "github.com/HardDie/DeckBuilder/internal/db/game"
	"github.com/HardDie/DeckBuilder/internal/entity"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	repositoriesCollection "github.com/HardDie/DeckBuilder/internal/repositories/collection"
	repositoriesDeck "github.com/HardDie/DeckBuilder/internal/repositories/deck"
	repositoriesGame "github.com/HardDie/DeckBuilder/internal/repositories/game"
	servicesCollection "github.com/HardDie/DeckBuilder/internal/services/collection"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type deckTest struct {
	gameID, collectionID string
	cfg                  *config.Config
	core                 dbCore.Core

	serviceGame       servicesGame.Game
	serviceCollection servicesCollection.Collection
	serviceDeck       Deck
}

func newDeckTest(t testing.TB) *deckTest {
	dir, err := os.MkdirTemp("", "deck_test")
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

	repositoryGame := repositoriesGame.New(cfg, game)
	repositoryCollection := repositoriesCollection.New(cfg, collection)
	repositoryDeck := repositoriesDeck.New(cfg, collection, deck)

	return &deckTest{
		gameID:       "test_deck__game",
		collectionID: "test_deck__collection",
		cfg:          cfg,
		core:         core,

		serviceGame:       servicesGame.New(cfg, repositoryGame),
		serviceCollection: servicesCollection.New(cfg, repositoryCollection),
		serviceDeck:       New(cfg, repositoryDeck),
	}
}

func (tt *deckTest) testCreate(t *testing.T) {
	deckType := "create_one"

	// Create deck
	deck, err := tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name: deckType,
	})
	if err != nil {
		t.Fatal(err)
	}
	if deck.Name != deckType {
		t.Fatal("Bad type [got]", deck.Name, "[want]", deckType)
	}

	// Try to create duplicate
	_, err = tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name: deckType,
	})
	if err == nil {
		t.Fatal("Error, you can't create duplicate deck")
	}
	if !errors.Is(err, er.DeckExist) {
		t.Fatal(err)
	}

	// Delete deck
	err = tt.serviceDeck.Delete(tt.gameID, tt.collectionID, deck.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *deckTest) testDelete(t *testing.T) {
	deckType := "delete_one"
	deckID := utils.NameToID(deckType)

	// Try to remove non-existing deck
	err := tt.serviceDeck.Delete(tt.gameID, tt.collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name: deckType,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete deck
	err = tt.serviceDeck.Delete(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete deck twice
	err = tt.serviceDeck.Delete(tt.gameID, tt.collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}
}
func (tt *deckTest) testUpdate(t *testing.T) {
	deckType := []string{"update_one", "update_two"}
	deckID := []string{utils.NameToID(deckType[0]), utils.NameToID(deckType[1])}

	// Try to update non-existing deck
	_, err := tt.serviceDeck.Update(tt.gameID, tt.collectionID, deckID[0], UpdateRequest{})
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	deck, err := tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name: deckType[0],
	})
	if err != nil {
		t.Fatal(err)
	}
	if deck.Name != deckType[0] {
		t.Fatal("Bad type [got]", deck.Name, "[want]", deckType[0])
	}

	// Update deck
	deck, err = tt.serviceDeck.Update(tt.gameID, tt.collectionID, deckID[0], UpdateRequest{
		Name: deckType[1],
	})
	if err != nil {
		t.Fatal(err)
	}
	if deck.Name != deckType[1] {
		t.Fatal("Bad type [got]", deck.Name, "[want]", deckType[1])
	}

	// Delete deck
	err = tt.serviceDeck.Delete(tt.gameID, tt.collectionID, deckID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Try to update non-existing deck
	_, err = tt.serviceDeck.Update(tt.gameID, tt.collectionID, deckID[1], UpdateRequest{})
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}
}
func (tt *deckTest) testList(t *testing.T) {
	deckType := []string{"B deck", "A deck"}
	deckID := []string{utils.NameToID(deckType[0]), utils.NameToID(deckType[1])}

	// Empty list
	items, _, err := tt.serviceDeck.List(tt.gameID, tt.collectionID, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}

	// Create first deck
	_, err = tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name: deckType[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// One deck
	items, _, err = tt.serviceDeck.List(tt.gameID, tt.collectionID, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatal("List should be with 1 element")
	}

	// Create second deck
	_, err = tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name: deckType[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sort by name
	items, _, err = tt.serviceDeck.List(tt.gameID, tt.collectionID, "name", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != deckType[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", deckType[1])
	}
	if items[1].Name != deckType[0] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", deckType[0])
	}

	// Sort by name_desc
	items, _, err = tt.serviceDeck.List(tt.gameID, tt.collectionID, "name_desc", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != deckType[0] {
		t.Fatal("Bad name_desc order: [got]", items[0].Name, "[want]", deckType[0])
	}
	if items[1].Name != deckType[1] {
		t.Fatal("Bad name_desc order: [got]", items[1].Name, "[want]", deckType[1])
	}

	// Sort by created date
	items, _, err = tt.serviceDeck.List(tt.gameID, tt.collectionID, "created", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != deckType[0] {
		t.Fatal("Bad created order: [got]", items[0].Name, "[want]", deckType[0])
	}
	if items[1].Name != deckType[1] {
		t.Fatal("Bad created order: [got]", items[1].Name, "[want]", deckType[1])
	}

	// Sort by created_desc
	items, _, err = tt.serviceDeck.List(tt.gameID, tt.collectionID, "created_desc", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != deckType[1] {
		t.Fatal("Bad created_desc order: [got]", items[0].Name, "[want]", deckType[1])
	}
	if items[1].Name != deckType[0] {
		t.Fatal("Bad created_desc order: [got]", items[1].Name, "[want]", deckType[0])
	}

	// Delete first deck
	err = tt.serviceDeck.Delete(tt.gameID, tt.collectionID, deckID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Delete second deck
	err = tt.serviceDeck.Delete(tt.gameID, tt.collectionID, deckID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Empty list
	items, _, err = tt.serviceDeck.List(tt.gameID, tt.collectionID, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}
}
func (tt *deckTest) testItem(t *testing.T) {
	deckType := []string{"item_one", "item_two"}
	deckID := []string{utils.NameToID(deckType[0]), utils.NameToID(deckType[1])}

	// Try to get non-existing deck
	_, err := tt.serviceDeck.Item(tt.gameID, tt.collectionID, deckID[0])
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name: deckType[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid deck
	_, err = tt.serviceDeck.Item(tt.gameID, tt.collectionID, deckID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid deck
	_, err = tt.serviceDeck.Item(tt.gameID, tt.collectionID, deckID[1])
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Rename deck
	_, err = tt.serviceDeck.Update(tt.gameID, tt.collectionID, deckID[0], UpdateRequest{Name: deckType[1]})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid deck
	_, err = tt.serviceDeck.Item(tt.gameID, tt.collectionID, deckID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid deck
	_, err = tt.serviceDeck.Item(tt.gameID, tt.collectionID, deckID[0])
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Delete deck
	err = tt.serviceDeck.Delete(tt.gameID, tt.collectionID, deckID[1])
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *deckTest) testImage(t *testing.T) {
	deckType := "image_one"
	deckID := utils.NameToID(deckType)
	pngImage := "https://github.com/fluidicon.png"
	jpegImage := "https://avatars.githubusercontent.com/apple"

	// Check no deck
	_, _, err := tt.serviceDeck.GetImage(tt.gameID, tt.collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck not exists")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name:  deckType,
		Image: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.serviceDeck.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update deck
	_, err = tt.serviceDeck.Update(tt.gameID, tt.collectionID, deckID, UpdateRequest{
		Name:  deckType,
		Image: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceDeck.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update deck
	_, err = tt.serviceDeck.Update(tt.gameID, tt.collectionID, deckID, UpdateRequest{
		Name:  deckType,
		Image: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.serviceDeck.GetImage(tt.gameID, tt.collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck don't have image")
	}
	if !errors.Is(err, er.DeckImageNotExists) {
		t.Fatal(err)
	}

	// Delete deck
	err = tt.serviceDeck.Delete(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *deckTest) testImageBin(t *testing.T) {
	deckType := "image_bin_one"
	deckID := utils.NameToID(deckType)

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

	// Check no deck
	_, _, err = tt.serviceDeck.GetImage(tt.gameID, tt.collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck not exists")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name:      deckType,
		ImageFile: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.serviceDeck.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update deck
	_, err = tt.serviceDeck.Update(tt.gameID, tt.collectionID, deckID, UpdateRequest{
		Name:      deckType,
		ImageFile: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceDeck.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update deck
	_, err = tt.serviceDeck.Update(tt.gameID, tt.collectionID, deckID, UpdateRequest{
		Name:      deckType,
		ImageFile: gifImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceDeck.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "gif" {
		t.Fatal("Image type error! [got]", imgType, "[want] gif")
	}

	// Update deck
	_, err = tt.serviceDeck.Update(tt.gameID, tt.collectionID, deckID, UpdateRequest{
		Name: deckType,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceDeck.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "gif" {
		t.Fatal("Image type error! [got]", imgType, "[want] gif")
	}

	// Update deck
	_, err = tt.serviceDeck.Update(tt.gameID, tt.collectionID, deckID, UpdateRequest{
		Name:  deckType,
		Image: "empty",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.serviceDeck.GetImage(tt.gameID, tt.collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck don't have image")
	}
	if !errors.Is(err, er.DeckImageNotExists) {
		t.Fatal(err)
	}

	// Delete deck
	err = tt.serviceDeck.Delete(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeck(t *testing.T) {
	t.Parallel()

	tt := newDeckTest(t)

	if err := tt.core.Init(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := tt.core.Drop(); err != nil {
			t.Fatal(err)
		}
	}()

	// Game not exist error
	_, err := tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name: "test",
	})
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = tt.serviceGame.Create(servicesGame.CreateRequest{
		Name: tt.gameID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Collection not exist error
	_, err = tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name: "test",
	})
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}

	// Create collection
	_, err = tt.serviceCollection.Create(tt.gameID, servicesCollection.CreateRequest{
		Name: tt.collectionID,
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
	t.Run("image_bin", tt.testImageBin)
}

func (tt *deckTest) fuzzCleanup() {
	_ = tt.core.Drop()
	_ = tt.core.Init()
}
func (tt *deckTest) fuzzList(t *testing.T, waitItems int) error {
	items, _, err := tt.serviceDeck.List(tt.gameID, tt.collectionID, "", "")
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
func (tt *deckTest) fuzzItem(t *testing.T, deckID, deckName string) error {
	deck, err := tt.serviceDeck.Item(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	if deck.Name != deckName {
		{
			data, _ := json.MarshalIndent(deck, "", "	")
			t.Log(string(data))
		}
		return fmt.Errorf("name: [wait] %s [got] %s", deckName, deck.Name)
	}
	return nil
}
func (tt *deckTest) fuzzCreate(t *testing.T, deckName string) (*entity.DeckInfo, error) {
	deck, err := tt.serviceDeck.Create(tt.gameID, tt.collectionID, CreateRequest{
		Name: deckName,
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
func (tt *deckTest) fuzzUpdate(t *testing.T, deckID, deckName string) (*entity.DeckInfo, error) {
	deck, err := tt.serviceDeck.Update(tt.gameID, tt.collectionID, deckID, UpdateRequest{
		Name: deckName,
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
func (tt *deckTest) fuzzDelete(t *testing.T, deckID string) error {
	err := tt.serviceDeck.Delete(tt.gameID, tt.collectionID, deckID)
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
	tt := newDeckTest(f)

	if err := tt.core.Init(); err != nil {
		f.Fatal(err)
	}
	defer func() {
		if err := tt.core.Drop(); err != nil {
			f.Fatal(err)
		}
	}()

	f.Fuzz(func(t *testing.T, type1, type2 string) {
		gameItems, _, err := tt.serviceGame.List("", "")
		if err != nil {
			t.Fatal(err)
		}
		if len(gameItems) == 0 {
			// Create game
			_, err = tt.serviceGame.Create(servicesGame.CreateRequest{
				Name: tt.gameID,
			})
			if err != nil {
				f.Fatal(err)
			}

			// Create collection
			_, err = tt.serviceCollection.Create(tt.gameID, servicesCollection.CreateRequest{
				Name: tt.collectionID,
			})
			if err != nil {
				f.Fatal(err)
			}
		}

		if utils.NameToID(type1) == "" || utils.NameToID(type2) == "" {
			// skip
			return
		}

		// Empty list
		err = tt.fuzzList(t, 0)
		if err != nil {
			t.Fatal(err)
		}

		// Create deck
		deck1, err := tt.fuzzCreate(t, type1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// List with deck
		err = tt.fuzzList(t, 1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = tt.fuzzItem(t, deck1.ID, type1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Update collection
		collection2, err := tt.fuzzUpdate(t, utils.NameToID(type1), type2)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// List with collection
		err = tt.fuzzList(t, 1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = tt.fuzzItem(t, collection2.ID, type2)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Delete collection
		err = tt.fuzzDelete(t, utils.NameToID(type2))
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
