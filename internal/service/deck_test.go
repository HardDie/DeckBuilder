package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/HardDie/fsentry"
	"github.com/google/uuid"

	"github.com/HardDie/DeckBuilder/internal/config"
	"github.com/HardDie/DeckBuilder/internal/db"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	"github.com/HardDie/DeckBuilder/internal/repository"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type deckTest struct {
	gameID, collectionID string
	cfg                  *config.Config
	gameService          IGameService
	collectionService    ICollectionService
	deckService          IDeckService
	db                   *db.DB
}

func newDeckTest(dataPath string) *deckTest {
	cfg := config.Get(false, "")
	cfg.SetDataPath(dataPath)

	// fsentry db
	builderDB := db.NewFSEntryDB(fsentry.NewFSEntry(cfg.Games()))

	gameRepository := repository.NewGameRepository(cfg, builderDB)
	collectionRepository := repository.NewCollectionRepository(cfg, builderDB)

	return &deckTest{
		gameID:            "test_deck__game",
		collectionID:      "test_deck__collection",
		cfg:               cfg,
		gameService:       NewGameService(cfg, gameRepository),
		collectionService: NewCollectionService(cfg, collectionRepository),
		deckService:       NewDeckService(cfg, repository.NewDeckRepository(cfg, builderDB)),
		db:                builderDB,
	}
}

func (tt *deckTest) testCreate(t *testing.T) {
	deckType := "create_one"

	// Create deck
	deck, err := tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
		Name: deckType,
	})
	if err != nil {
		t.Fatal(err)
	}
	if deck.Name != deckType {
		t.Fatal("Bad type [got]", deck.Name, "[want]", deckType)
	}

	// Try to create duplicate
	_, err = tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
		Name: deckType,
	})
	if err == nil {
		t.Fatal("Error, you can't create duplicate deck")
	}
	if !errors.Is(err, er.DeckExist) {
		t.Fatal(err)
	}

	// Delete deck
	err = tt.deckService.Delete(tt.gameID, tt.collectionID, deck.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *deckTest) testDelete(t *testing.T) {
	deckType := "delete_one"
	deckID := utils.NameToID(deckType)

	// Try to remove non-existing deck
	err := tt.deckService.Delete(tt.gameID, tt.collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
		Name: deckType,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete deck
	err = tt.deckService.Delete(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete deck twice
	err = tt.deckService.Delete(tt.gameID, tt.collectionID, deckID)
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
	_, err := tt.deckService.Update(tt.gameID, tt.collectionID, deckID[0], &dto.UpdateDeckDTO{})
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	deck, err := tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
		Name: deckType[0],
	})
	if err != nil {
		t.Fatal(err)
	}
	if deck.Name != deckType[0] {
		t.Fatal("Bad type [got]", deck.Name, "[want]", deckType[0])
	}

	// Update deck
	deck, err = tt.deckService.Update(tt.gameID, tt.collectionID, deckID[0], &dto.UpdateDeckDTO{
		Name: deckType[1],
	})
	if err != nil {
		t.Fatal(err)
	}
	if deck.Name != deckType[1] {
		t.Fatal("Bad type [got]", deck.Name, "[want]", deckType[1])
	}

	// Delete deck
	err = tt.deckService.Delete(tt.gameID, tt.collectionID, deckID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Try to update non-existing deck
	_, err = tt.deckService.Update(tt.gameID, tt.collectionID, deckID[1], &dto.UpdateDeckDTO{})
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
	items, err := tt.deckService.List(tt.gameID, tt.collectionID, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}

	// Create first deck
	_, err = tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
		Name: deckType[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// One deck
	items, err = tt.deckService.List(tt.gameID, tt.collectionID, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatal("List should be with 1 element")
	}

	// Create second deck
	_, err = tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
		Name: deckType[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sort by name
	items, err = tt.deckService.List(tt.gameID, tt.collectionID, "name", "")
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
	items, err = tt.deckService.List(tt.gameID, tt.collectionID, "name_desc", "")
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
	items, err = tt.deckService.List(tt.gameID, tt.collectionID, "created", "")
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
	items, err = tt.deckService.List(tt.gameID, tt.collectionID, "created_desc", "")
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
	err = tt.deckService.Delete(tt.gameID, tt.collectionID, deckID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Delete second deck
	err = tt.deckService.Delete(tt.gameID, tt.collectionID, deckID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Empty list
	items, err = tt.deckService.List(tt.gameID, tt.collectionID, "", "")
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
	_, err := tt.deckService.Item(tt.gameID, tt.collectionID, deckID[0])
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
		Name: deckType[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid deck
	_, err = tt.deckService.Item(tt.gameID, tt.collectionID, deckID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid deck
	_, err = tt.deckService.Item(tt.gameID, tt.collectionID, deckID[1])
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Rename deck
	_, err = tt.deckService.Update(tt.gameID, tt.collectionID, deckID[0], &dto.UpdateDeckDTO{Name: deckType[1]})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid deck
	_, err = tt.deckService.Item(tt.gameID, tt.collectionID, deckID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid deck
	_, err = tt.deckService.Item(tt.gameID, tt.collectionID, deckID[0])
	if err == nil {
		t.Fatal("Error, deck not exist")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Delete deck
	err = tt.deckService.Delete(tt.gameID, tt.collectionID, deckID[1])
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
	_, _, err := tt.deckService.GetImage(tt.gameID, tt.collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck not exists")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
		Name:  deckType,
		Image: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.deckService.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update deck
	_, err = tt.deckService.Update(tt.gameID, tt.collectionID, deckID, &dto.UpdateDeckDTO{
		Name:  deckType,
		Image: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.deckService.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update deck
	_, err = tt.deckService.Update(tt.gameID, tt.collectionID, deckID, &dto.UpdateDeckDTO{
		Name:  deckType,
		Image: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.deckService.GetImage(tt.gameID, tt.collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck don't have image")
	}
	if !errors.Is(err, er.DeckImageNotExists) {
		t.Fatal(err)
	}

	// Delete deck
	err = tt.deckService.Delete(tt.gameID, tt.collectionID, deckID)
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
	_, _, err = tt.deckService.GetImage(tt.gameID, tt.collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck not exists")
	}
	if !errors.Is(err, er.DeckNotExists) {
		t.Fatal(err)
	}

	// Create deck
	_, err = tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
		Name:      deckType,
		ImageFile: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.deckService.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update deck
	_, err = tt.deckService.Update(tt.gameID, tt.collectionID, deckID, &dto.UpdateDeckDTO{
		Name:      deckType,
		ImageFile: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.deckService.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update deck
	_, err = tt.deckService.Update(tt.gameID, tt.collectionID, deckID, &dto.UpdateDeckDTO{
		Name:      deckType,
		ImageFile: gifImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.deckService.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "gif" {
		t.Fatal("Image type error! [got]", imgType, "[want] gif")
	}

	// Update deck
	_, err = tt.deckService.Update(tt.gameID, tt.collectionID, deckID, &dto.UpdateDeckDTO{
		Name: deckType,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.deckService.GetImage(tt.gameID, tt.collectionID, deckID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "gif" {
		t.Fatal("Image type error! [got]", imgType, "[want] gif")
	}

	// Update deck
	_, err = tt.deckService.Update(tt.gameID, tt.collectionID, deckID, &dto.UpdateDeckDTO{
		Name:  deckType,
		Image: "empty",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.deckService.GetImage(tt.gameID, tt.collectionID, deckID)
	if err == nil {
		t.Fatal("Error, deck don't have image")
	}
	if !errors.Is(err, er.DeckImageNotExists) {
		t.Fatal(err)
	}

	// Delete deck
	err = tt.deckService.Delete(tt.gameID, tt.collectionID, deckID)
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
	tt := newDeckTest(filepath.Join(dataPath, "deck_test"))

	if err := tt.db.Init(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := tt.db.Drop(); err != nil {
			t.Fatal(err)
		}
	}()

	// Game not exist error
	_, err := tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
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
	_, err = tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
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

	t.Run("create", tt.testCreate)
	t.Run("delete", tt.testDelete)
	t.Run("update", tt.testUpdate)
	t.Run("list", tt.testList)
	t.Run("item", tt.testItem)
	t.Run("image", tt.testImage)
	t.Run("image_bin", tt.testImageBin)
}

func (tt *deckTest) fuzzCleanup() {
	_ = tt.db.Drop()
	_ = tt.db.Init()
}
func (tt *deckTest) fuzzList(t *testing.T, waitItems int) error {
	items, err := tt.deckService.List(tt.gameID, tt.collectionID, "", "")
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
	deck, err := tt.deckService.Item(tt.gameID, tt.collectionID, deckID)
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
	deck, err := tt.deckService.Create(tt.gameID, tt.collectionID, &dto.CreateDeckDTO{
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
	deck, err := tt.deckService.Update(tt.gameID, tt.collectionID, deckID, &dto.UpdateDeckDTO{
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
	err := tt.deckService.Delete(tt.gameID, tt.collectionID, deckID)
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
	tt := newDeckTest(filepath.Join(dataPath, "deck_fuzz_"+uuid.New().String()))

	if err := tt.db.Init(); err != nil {
		f.Fatal(err)
	}
	defer func() {
		if err := tt.db.Drop(); err != nil {
			f.Fatal(err)
		}
	}()

	f.Fuzz(func(t *testing.T, type1, type2 string) {
		gameItems, err := tt.gameService.List("", "")
		if err != nil {
			t.Fatal(err)
		}
		if len(gameItems) == 0 {
			// Create game
			_, err = tt.gameService.Create(&dto.CreateGameDTO{
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
