package collection

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
	dbGame "github.com/HardDie/DeckBuilder/internal/db/game"
	entitiesCollection "github.com/HardDie/DeckBuilder/internal/entities/collection"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	repositoriesCollection "github.com/HardDie/DeckBuilder/internal/repositories/collection"
	repositoriesGame "github.com/HardDie/DeckBuilder/internal/repositories/game"
	servicesGame "github.com/HardDie/DeckBuilder/internal/services/game"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type collectionTest struct {
	gameID string
	cfg    *config.Config
	core   dbCore.Core

	serviceGame       servicesGame.Game
	serviceCollection Collection
}

func newCollectionTest(t testing.TB) *collectionTest {
	dir, err := os.MkdirTemp("", "collection_test")
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

	repositoryGame := repositoriesGame.New(cfg, game)
	repositoryCollection := repositoriesCollection.New(cfg, collection)

	return &collectionTest{
		gameID: "test_collection__game",
		cfg:    cfg,
		core:   core,

		serviceGame:       servicesGame.New(cfg, repositoryGame),
		serviceCollection: New(cfg, repositoryCollection),
	}
}

func (tt *collectionTest) testCreate(t *testing.T) {
	collectionName := "create_one"
	desc := "best game ever"

	// Create collection
	collection, err := tt.serviceCollection.Create(tt.gameID, CreateRequest{
		Name:        collectionName,
		Description: desc,
	})
	if err != nil {
		t.Fatal(err)
	}
	if collection.Name != collectionName {
		t.Fatal("Bad name [got]", collection.Name, "[want]", collectionName)
	}
	if collection.Description != desc {
		t.Fatal("Bad description [got]", collection.Description, "[want]", desc)
	}

	// Try to create duplicate
	_, err = tt.serviceCollection.Create(tt.gameID, CreateRequest{
		Name: collectionName,
	})
	if err == nil {
		t.Fatal("Error, you can't create duplicate collection")
	}
	if !errors.Is(err, er.CollectionExist) {
		t.Fatal(err)
	}

	// Delete collection
	err = tt.serviceCollection.Delete(tt.gameID, collection.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *collectionTest) testDelete(t *testing.T) {
	collectionName := "delete_one"
	collectionID := utils.NameToID(collectionName)

	// Try to remove non-existing collection
	err := tt.serviceCollection.Delete(tt.gameID, collectionID)
	if err == nil {
		t.Fatal("Error, collection not exist")
	}
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}

	// Create collection
	_, err = tt.serviceCollection.Create(tt.gameID, CreateRequest{
		Name: collectionName,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete collection
	err = tt.serviceCollection.Delete(tt.gameID, collectionID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete collection twice
	err = tt.serviceCollection.Delete(tt.gameID, collectionID)
	if err == nil {
		t.Fatal("Error, collection not exist")
	}
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}
}
func (tt *collectionTest) testUpdate(t *testing.T) {
	collectionName := []string{"update_one", "update_two"}
	desc := []string{"first description", "second description"}
	collectionID := []string{utils.NameToID(collectionName[0]), utils.NameToID(collectionName[1])}

	// Try to update non-existing collection
	_, err := tt.serviceCollection.Update(tt.gameID, collectionID[0], UpdateRequest{})
	if err == nil {
		t.Fatal("Error, collection not exist")
	}
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}

	// Create collection
	collection, err := tt.serviceCollection.Create(tt.gameID, CreateRequest{
		Name:        collectionName[0],
		Description: desc[0],
	})
	if err != nil {
		t.Fatal(err)
	}
	if collection.Name != collectionName[0] {
		t.Fatal("Bad name [got]", collection.Name, "[want]", collectionName[0])
	}
	if collection.Description != desc[0] {
		t.Fatal("Bad description [got]", collection.Description, "[want]", desc[0])
	}

	// Update collection
	collection, err = tt.serviceCollection.Update(tt.gameID, collectionID[0], UpdateRequest{
		Name:        collectionName[1],
		Description: desc[1],
	})
	if err != nil {
		t.Fatal(err)
	}
	if collection.Name != collectionName[1] {
		t.Fatal("Bad name [got]", collection.Name, "[want]", collectionName[1])
	}
	if collection.Description != desc[1] {
		t.Fatal("Bad description [got]", collection.Description, "[want]", desc[1])
	}

	// Delete collection
	err = tt.serviceCollection.Delete(tt.gameID, collectionID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Try to update non-existing collection
	_, err = tt.serviceCollection.Update(tt.gameID, collectionID[1], UpdateRequest{})
	if err == nil {
		t.Fatal("Error, collection not exist")
	}
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}
}
func (tt *collectionTest) testList(t *testing.T) {
	collectionName := []string{"B collection", "A collection"}
	collectionID := []string{utils.NameToID(collectionName[0]), utils.NameToID(collectionName[1])}

	// Empty list
	items, err := tt.serviceCollection.List(tt.gameID, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}

	// Create first collection
	_, err = tt.serviceCollection.Create(tt.gameID, CreateRequest{
		Name: collectionName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// One collection
	items, err = tt.serviceCollection.List(tt.gameID, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatal("List should be with 1 element")
	}

	// Create second collection
	_, err = tt.serviceCollection.Create(tt.gameID, CreateRequest{
		Name: collectionName[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sort by name
	items, err = tt.serviceCollection.List(tt.gameID, "name", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != collectionName[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", collectionName[1])
	}
	if items[1].Name != collectionName[0] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", collectionName[0])
	}

	// Sort by name_desc
	items, err = tt.serviceCollection.List(tt.gameID, "name_desc", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != collectionName[0] {
		t.Fatal("Bad name_desc order: [got]", items[0].Name, "[want]", collectionName[0])
	}
	if items[1].Name != collectionName[1] {
		t.Fatal("Bad name_desc order: [got]", items[1].Name, "[want]", collectionName[1])
	}

	// Sort by created date
	items, err = tt.serviceCollection.List(tt.gameID, "created", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != collectionName[0] {
		t.Fatal("Bad created order: [got]", items[0].Name, "[want]", collectionName[0])
	}
	if items[1].Name != collectionName[1] {
		t.Fatal("Bad created order: [got]", items[1].Name, "[want]", collectionName[1])
	}

	// Sort by created_desc
	items, err = tt.serviceCollection.List(tt.gameID, "created_desc", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != collectionName[1] {
		t.Fatal("Bad created_desc order: [got]", items[0].Name, "[want]", collectionName[1])
	}
	if items[1].Name != collectionName[0] {
		t.Fatal("Bad created_desc order: [got]", items[1].Name, "[want]", collectionName[0])
	}

	// Delete first collection
	err = tt.serviceCollection.Delete(tt.gameID, collectionID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Delete second collection
	err = tt.serviceCollection.Delete(tt.gameID, collectionID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Empty list
	items, err = tt.serviceCollection.List(tt.gameID, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}
}
func (tt *collectionTest) testItem(t *testing.T) {
	collectionName := []string{"item_one", "item_two"}
	collectionID := []string{utils.NameToID(collectionName[0]), utils.NameToID(collectionName[1])}

	// Try to get non-existing collection
	_, err := tt.serviceCollection.Item(tt.gameID, collectionID[0])
	if err == nil {
		t.Fatal("Error, collection not exist")
	}
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}

	// Create collection
	_, err = tt.serviceCollection.Create(tt.gameID, CreateRequest{
		Name: collectionName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid collection
	_, err = tt.serviceCollection.Item(tt.gameID, collectionID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid collection
	_, err = tt.serviceCollection.Item(tt.gameID, collectionID[1])
	if err == nil {
		t.Fatal("Error, collection not exist")
	}
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}

	// Rename collection
	_, err = tt.serviceCollection.Update(tt.gameID, collectionID[0], UpdateRequest{Name: collectionName[1]})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid collection
	_, err = tt.serviceCollection.Item(tt.gameID, collectionID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid collection
	_, err = tt.serviceCollection.Item(tt.gameID, collectionID[0])
	if err == nil {
		t.Fatal("Error, collection not exist")
	}
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}

	// Delete collection
	err = tt.serviceCollection.Delete(tt.gameID, collectionID[1])
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *collectionTest) testImage(t *testing.T) {
	collectionName := "image_one"
	collectionID := utils.NameToID(collectionName)
	pngImage := "https://github.com/fluidicon.png"
	jpegImage := "https://avatars.githubusercontent.com/apple"

	// Check no collection
	_, _, err := tt.serviceCollection.GetImage(tt.gameID, collectionID)
	if err == nil {
		t.Fatal("Error, collection not exists")
	}
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}

	// Create collection
	_, err = tt.serviceCollection.Create(tt.gameID, CreateRequest{
		Name:  collectionName,
		Image: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.serviceCollection.GetImage(tt.gameID, collectionID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update collection
	_, err = tt.serviceCollection.Update(tt.gameID, collectionID, UpdateRequest{
		Name:  collectionName,
		Image: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceCollection.GetImage(tt.gameID, collectionID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update collection
	_, err = tt.serviceCollection.Update(tt.gameID, collectionID, UpdateRequest{
		Name:  collectionName,
		Image: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.serviceCollection.GetImage(tt.gameID, collectionID)
	if err == nil {
		t.Fatal("Error, collection don't have image")
	}
	if !errors.Is(err, er.CollectionImageNotExists) {
		t.Fatal(err)
	}

	// Delete collection
	err = tt.serviceCollection.Delete(tt.gameID, collectionID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *collectionTest) testImageBin(t *testing.T) {
	collectionName := "image_bin_one"
	collectionID := utils.NameToID(collectionName)

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

	// Check no collection
	_, _, err = tt.serviceCollection.GetImage(tt.gameID, collectionID)
	if err == nil {
		t.Fatal("Error, collection not exists")
	}
	if !errors.Is(err, er.CollectionNotExists) {
		t.Fatal(err)
	}

	// Create collection
	_, err = tt.serviceCollection.Create(tt.gameID, CreateRequest{
		Name:      collectionName,
		ImageFile: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.serviceCollection.GetImage(tt.gameID, collectionID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update collection
	_, err = tt.serviceCollection.Update(tt.gameID, collectionID, UpdateRequest{
		Name:      collectionName,
		ImageFile: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceCollection.GetImage(tt.gameID, collectionID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update collection
	_, err = tt.serviceCollection.Update(tt.gameID, collectionID, UpdateRequest{
		Name:      collectionName,
		ImageFile: gifImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceCollection.GetImage(tt.gameID, collectionID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "gif" {
		t.Fatal("Image type error! [got]", imgType, "[want] gif")
	}

	// Update collection
	_, err = tt.serviceCollection.Update(tt.gameID, collectionID, UpdateRequest{
		Name: collectionName,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceCollection.GetImage(tt.gameID, collectionID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "gif" {
		t.Fatal("Image type error! [got]", imgType, "[want] gif")
	}

	// Update collection
	_, err = tt.serviceCollection.Update(tt.gameID, collectionID, UpdateRequest{
		Name:  collectionName,
		Image: "empty",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.serviceCollection.GetImage(tt.gameID, collectionID)
	if err == nil {
		t.Fatal("Error, collection don't have image")
	}
	if !errors.Is(err, er.CollectionImageNotExists) {
		t.Fatal(err)
	}

	// Delete collection
	err = tt.serviceCollection.Delete(tt.gameID, collectionID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCollection(t *testing.T) {
	t.Parallel()

	tt := newCollectionTest(t)

	if err := tt.core.Init(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := tt.core.Drop(); err != nil {
			t.Fatal(err)
		}
	}()

	// Game not exist error
	_, err := tt.serviceCollection.Create(tt.gameID, CreateRequest{
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

	t.Run("create", tt.testCreate)
	t.Run("delete", tt.testDelete)
	t.Run("update", tt.testUpdate)
	t.Run("list", tt.testList)
	t.Run("item", tt.testItem)
	t.Run("image", tt.testImage)
	t.Run("image_bin", tt.testImageBin)
}

func (tt *collectionTest) fuzzCleanup() {
	_ = tt.core.Drop()
	_ = tt.core.Init()
}
func (tt *collectionTest) fuzzList(t *testing.T, waitItems int) error {
	items, err := tt.serviceCollection.List(tt.gameID, "", "")
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
func (tt *collectionTest) fuzzItem(t *testing.T, collectionID, name, desc string) error {
	collection, err := tt.serviceCollection.Item(tt.gameID, collectionID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	if collection.Name != name {
		{
			data, _ := json.MarshalIndent(collection, "", "	")
			t.Log(string(data))
		}
		return fmt.Errorf("name: [wait] %s [got] %s", name, collection.Name)
	}
	if collection.Description != desc {
		{
			data, _ := json.MarshalIndent(collection, "", "	")
			t.Log("item:", string(data))
		}
		return fmt.Errorf("description: [wait] %q [got] %q", desc, collection.Description)
	}
	return nil
}
func (tt *collectionTest) fuzzCreate(t *testing.T, name, desc string) (*entitiesCollection.Collection, error) {
	collection, err := tt.serviceCollection.Create(tt.gameID, CreateRequest{
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
		data, _ := json.MarshalIndent(collection, "", "	")
		t.Log("create:", string(data))
	}
	return collection, nil
}
func (tt *collectionTest) fuzzUpdate(t *testing.T, collectionID, name, desc string) (*entitiesCollection.Collection, error) {
	collection, err := tt.serviceCollection.Update(tt.gameID, collectionID, UpdateRequest{
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
		data, _ := json.MarshalIndent(collection, "", "	")
		t.Log("update:", string(data))
	}
	return collection, nil
}
func (tt *collectionTest) fuzzDelete(t *testing.T, collectionID string) error {
	err := tt.serviceCollection.Delete(tt.gameID, collectionID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	return nil
}

func FuzzCollection(f *testing.F) {
	tt := newCollectionTest(f)

	if err := tt.core.Init(); err != nil {
		f.Fatal(err)
	}
	defer func() {
		if err := tt.core.Drop(); err != nil {
			f.Fatal(err)
		}
	}()

	f.Fuzz(func(t *testing.T, name1, desc1, name2, desc2 string) {
		items, err := tt.serviceGame.List("", "")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) == 0 {
			// Create game
			_, err := tt.serviceGame.Create(servicesGame.CreateRequest{
				Name: tt.gameID,
			})
			if err != nil {
				f.Fatal(err)
			}

		}

		if utils.NameToID(name1) == "" || utils.NameToID(name2) == "" {
			// skip
			return
		}

		// Empty list
		err = tt.fuzzList(t, 0)
		if err != nil {
			t.Fatal(err)
		}

		// Create collection
		collection1, err := tt.fuzzCreate(t, name1, desc1)
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
		err = tt.fuzzItem(t, collection1.ID, name1, desc1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Update collection
		collection2, err := tt.fuzzUpdate(t, utils.NameToID(name1), name2, desc2)
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
		err = tt.fuzzItem(t, collection2.ID, name2, desc2)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Delete collection
		err = tt.fuzzDelete(t, utils.NameToID(name2))
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
