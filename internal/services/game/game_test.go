package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/HardDie/fsentry"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCore "github.com/HardDie/DeckBuilder/internal/db/core"
	dbGame "github.com/HardDie/DeckBuilder/internal/db/game"
	"github.com/HardDie/DeckBuilder/internal/dto"
	"github.com/HardDie/DeckBuilder/internal/entity"
	er "github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/images"
	repositoriesGame "github.com/HardDie/DeckBuilder/internal/repositories/game"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type gameTest struct {
	cfg  *config.Config
	core dbCore.Core

	serviceGame Game
}

func newGameTest(t testing.TB) *gameTest {
	dir, err := os.MkdirTemp("", "game_test")
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

	repositoryGame := repositoriesGame.New(cfg, game)

	return &gameTest{
		cfg:  cfg,
		core: core,

		serviceGame: New(cfg, repositoryGame),
	}
}

func (tt *gameTest) testCreate(t *testing.T) {
	gameName := "create_one"
	desc := "best game ever"

	// Create game
	game, err := tt.serviceGame.Create(&dto.CreateGameDTO{
		Name:        gameName,
		Description: desc,
	})
	if err != nil {
		t.Fatal(err)
	}
	if game.Name != gameName {
		t.Fatal("Bad name [got]", game.Name, "[want]", gameName)
	}
	if game.Description != desc {
		t.Fatal("Bad description [got]", game.Description, "[want]", desc)
	}

	// Try to create duplicate
	_, err = tt.serviceGame.Create(&dto.CreateGameDTO{
		Name: gameName,
	})
	if err == nil {
		t.Fatal("Error, you can't create duplicate game")
	}
	if !errors.Is(err, er.GameExist) {
		t.Fatal(err)
	}

	// Delete game
	err = tt.serviceGame.Delete(game.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func (tt *gameTest) testDelete(t *testing.T) {
	gameName := "delete_one"
	gameID := utils.NameToID(gameName)

	// Try to remove non-existing game
	err := tt.serviceGame.Delete(gameID)
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = tt.serviceGame.Create(&dto.CreateGameDTO{
		Name: gameName,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete game
	err = tt.serviceGame.Delete(gameID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete game twice
	err = tt.serviceGame.Delete(gameID)
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}
}
func (tt *gameTest) testUpdate(t *testing.T) {
	gameName := []string{"update_one", "update_two"}
	desc := []string{"first description", "second description"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Try to update non-existing game
	_, err := tt.serviceGame.Update(gameID[0], &dto.UpdateGameDTO{})
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	game, err := tt.serviceGame.Create(&dto.CreateGameDTO{
		Name:        gameName[0],
		Description: desc[0],
	})
	if err != nil {
		t.Fatal(err)
	}
	if game.Name != gameName[0] {
		t.Fatal("Bad name [got]", game.Name, "[want]", gameName[0])
	}
	if game.Description != desc[0] {
		t.Fatal("Bad description [got]", game.Description, "[want]", desc[0])
	}

	// Update game
	game, err = tt.serviceGame.Update(gameID[0], &dto.UpdateGameDTO{
		Name:        gameName[1],
		Description: desc[1],
	})
	if err != nil {
		t.Fatal(err)
	}
	if game.Name != gameName[1] {
		t.Fatal("Bad name [got]", game.Name, "[want]", gameName[1])
	}
	if game.Description != desc[1] {
		t.Fatal("Bad description [got]", game.Description, "[want]", desc[1])
	}

	// Delete game
	err = tt.serviceGame.Delete(gameID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Try to update non-existing game
	_, err = tt.serviceGame.Update(gameID[1], &dto.UpdateGameDTO{})
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}
}
func (tt *gameTest) testList(t *testing.T) {
	gameName := []string{"B game", "A game"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Empty list
	items, _, err := tt.serviceGame.List("", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}

	// Create first game
	_, err = tt.serviceGame.Create(&dto.CreateGameDTO{
		Name: gameName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// One game
	items, _, err = tt.serviceGame.List("", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatal("List should be with 1 element")
	}

	// Create second game
	_, err = tt.serviceGame.Create(&dto.CreateGameDTO{
		Name: gameName[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sort by name
	items, _, err = tt.serviceGame.List("name", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != gameName[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[1])
	}
	if items[1].Name != gameName[0] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[0])
	}

	// Sort by name_desc
	items, _, err = tt.serviceGame.List("name_desc", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != gameName[0] {
		t.Fatal("Bad name_desc order: [got]", items[0].Name, "[want]", gameName[0])
	}
	if items[1].Name != gameName[1] {
		t.Fatal("Bad name_desc order: [got]", items[1].Name, "[want]", gameName[1])
	}

	// Sort by created date
	items, _, err = tt.serviceGame.List("created", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != gameName[0] {
		t.Fatal("Bad created order: [got]", items[0].Name, "[want]", gameName[0])
	}
	if items[1].Name != gameName[1] {
		t.Fatal("Bad created order: [got]", items[1].Name, "[want]", gameName[1])
	}

	// Sort by created_desc
	items, _, err = tt.serviceGame.List("created_desc", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name != gameName[1] {
		t.Fatal("Bad created_desc order: [got]", items[0].Name, "[want]", gameName[1])
	}
	if items[1].Name != gameName[0] {
		t.Fatal("Bad created_desc order: [got]", items[1].Name, "[want]", gameName[0])
	}

	// Delete first game
	err = tt.serviceGame.Delete(gameID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Delete second game
	err = tt.serviceGame.Delete(gameID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Empty list
	items, _, err = tt.serviceGame.List("", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}
}
func (tt *gameTest) testItem(t *testing.T) {
	gameName := []string{"item_one", "item_two"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Try to get non-existing game
	_, err := tt.serviceGame.Item(gameID[0])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = tt.serviceGame.Create(&dto.CreateGameDTO{
		Name: gameName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid game
	_, err = tt.serviceGame.Item(gameID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid game
	_, err = tt.serviceGame.Item(gameID[1])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Rename game
	_, err = tt.serviceGame.Update(gameID[0], &dto.UpdateGameDTO{Name: gameName[1]})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid game
	_, err = tt.serviceGame.Item(gameID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid game
	_, err = tt.serviceGame.Item(gameID[0])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Delete game
	err = tt.serviceGame.Delete(gameID[1])
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *gameTest) testDuplicate(t *testing.T) {
	gameName := []string{"duplicate_one", "duplicate_two"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Create games
	_, err := tt.serviceGame.Create(&dto.CreateGameDTO{
		Name: gameName[0],
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = tt.serviceGame.Create(&dto.CreateGameDTO{
		Name: gameName[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Try to duplicate not exist game
	_, err = tt.serviceGame.Duplicate("not_exist_game", &dto.DuplicateGameDTO{
		Name: "new_game",
	})
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal("Game not exist")
	}

	// Try to duplicate to exist game
	_, err = tt.serviceGame.Duplicate(gameID[0], &dto.DuplicateGameDTO{
		Name: gameID[1],
	})
	if !errors.Is(err, er.GameExist) {
		t.Fatal("Game already exist")
	}

	_, err = tt.serviceGame.Duplicate(gameID[0], &dto.DuplicateGameDTO{
		Name: "good_duplicate",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = tt.serviceGame.Delete(gameID[0])
	if err != nil {
		t.Fatal(err)
	}
	err = tt.serviceGame.Delete(gameID[1])
	if err != nil {
		t.Fatal(err)
	}
	err = tt.serviceGame.Delete("good_duplicate")
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *gameTest) testImage(t *testing.T) {
	gameName := "image_one"
	gameID := utils.NameToID(gameName)
	pngImage := "https://github.com/fluidicon.png"
	jpegImage := "https://avatars.githubusercontent.com/apple"

	// Check no game
	_, _, err := tt.serviceGame.GetImage(gameID)
	if err == nil {
		t.Fatal("Error, game not exists")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = tt.serviceGame.Create(&dto.CreateGameDTO{
		Name:  gameName,
		Image: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.serviceGame.GetImage(gameID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update game
	_, err = tt.serviceGame.Update(gameID, &dto.UpdateGameDTO{
		Name:  gameName,
		Image: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceGame.GetImage(gameID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update game
	_, err = tt.serviceGame.Update(gameID, &dto.UpdateGameDTO{
		Name:  gameName,
		Image: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.serviceGame.GetImage(gameID)
	if err == nil {
		t.Fatal("Error, game don't have image")
	}
	if !errors.Is(err, er.GameImageNotExists) {
		t.Fatal(err)
	}

	// Delete game
	err = tt.serviceGame.Delete(gameID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *gameTest) testImageBin(t *testing.T) {
	gameName := "image_bin_one"
	gameID := utils.NameToID(gameName)

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

	// Check no game
	_, _, err = tt.serviceGame.GetImage(gameID)
	if err == nil {
		t.Fatal("Error, game not exists")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = tt.serviceGame.Create(&dto.CreateGameDTO{
		Name:      gameName,
		ImageFile: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.serviceGame.GetImage(gameID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update game
	_, err = tt.serviceGame.Update(gameID, &dto.UpdateGameDTO{
		Name:      gameName,
		ImageFile: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceGame.GetImage(gameID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update game
	_, err = tt.serviceGame.Update(gameID, &dto.UpdateGameDTO{
		Name:      gameName,
		ImageFile: gifImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceGame.GetImage(gameID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "gif" {
		t.Fatal("Image type error! [got]", imgType, "[want] gif")
	}

	// Update game
	_, err = tt.serviceGame.Update(gameID, &dto.UpdateGameDTO{
		Name: gameName,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.serviceGame.GetImage(gameID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "gif" {
		t.Fatal("Image type error! [got]", imgType, "[want] gif")
	}

	// Update game
	_, err = tt.serviceGame.Update(gameID, &dto.UpdateGameDTO{
		Name:  gameName,
		Image: "empty",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.serviceGame.GetImage(gameID)
	if err == nil {
		t.Fatal("Error, game don't have image")
	}
	if !errors.Is(err, er.GameImageNotExists) {
		t.Fatal(err)
	}

	// Delete game
	err = tt.serviceGame.Delete(gameID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGame(t *testing.T) {
	t.Parallel()

	tt := newGameTest(t)

	if err := tt.core.Init(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := tt.core.Drop(); err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("create", tt.testCreate)
	t.Run("delete", tt.testDelete)
	t.Run("update", tt.testUpdate)
	t.Run("list", tt.testList)
	t.Run("item", tt.testItem)
	t.Run("duplicate", tt.testDuplicate)
	t.Run("image", tt.testImage)
	t.Run("image_bin", tt.testImageBin)
}

func (tt *gameTest) fuzzCleanup() {
	_ = tt.core.Drop()
	_ = tt.core.Init()
}
func (tt *gameTest) fuzzList(t *testing.T, waitItems int) error {
	items, _, err := tt.serviceGame.List("", "")
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
func (tt *gameTest) fuzzItem(t *testing.T, gameID, name, desc string) error {
	game, err := tt.serviceGame.Item(gameID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	if game.Name != name {
		{
			data, _ := json.MarshalIndent(game, "", "	")
			t.Log(string(data))
		}
		return fmt.Errorf("name: [wait] %s [got] %s", name, game.Name)
	}
	if game.Description != desc {
		{
			data, _ := json.MarshalIndent(game, "", "	")
			t.Log("item:", string(data))
		}
		return fmt.Errorf("description: [wait] %q [got] %q", desc, game.Description)
	}
	return nil
}
func (tt *gameTest) fuzzCreate(t *testing.T, name, desc string) (*entity.GameInfo, error) {
	game, err := tt.serviceGame.Create(&dto.CreateGameDTO{
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
		data, _ := json.MarshalIndent(game, "", "	")
		t.Log("create:", string(data))
	}
	return game, nil
}
func (tt *gameTest) fuzzUpdate(t *testing.T, gameID, name, desc string) (*entity.GameInfo, error) {
	game, err := tt.serviceGame.Update(gameID, &dto.UpdateGameDTO{
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
		data, _ := json.MarshalIndent(game, "", "	")
		t.Log("update:", string(data))
	}
	return game, nil
}
func (tt *gameTest) fuzzDelete(t *testing.T, gameID string) error {
	err := tt.serviceGame.Delete(gameID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	return nil
}

func FuzzGame(f *testing.F) {
	tt := newGameTest(f)

	if err := tt.core.Init(); err != nil {
		f.Fatal(err)
	}
	defer func() {
		if err := tt.core.Drop(); err != nil {
			f.Fatal(err)
		}
	}()

	f.Fuzz(func(t *testing.T, name1, desc1, name2, desc2 string) {
		if utils.NameToID(name1) == "" || utils.NameToID(name2) == "" {
			// skip
			return
		}

		// Empty list
		err := tt.fuzzList(t, 0)
		if err != nil {
			t.Fatal(err)
		}

		// Create game
		game1, err := tt.fuzzCreate(t, name1, desc1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// List with game
		err = tt.fuzzList(t, 1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = tt.fuzzItem(t, game1.ID, name1, desc1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Update game
		game2, err := tt.fuzzUpdate(t, utils.NameToID(name1), name2, desc2)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// List with game
		err = tt.fuzzList(t, 1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = tt.fuzzItem(t, game2.ID, name2, desc2)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Delete game
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
