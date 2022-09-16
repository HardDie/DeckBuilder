package games

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/google/uuid"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	er "tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

func testCreate(t *testing.T) {
	service := NewService()
	gameName := "one"
	desc := "best game ever"

	// Create game
	game, err := service.Create(&dto.CreateGameDTO{
		Name:        gameName,
		Description: desc,
	})
	if err != nil {
		t.Fatal(err)
	}
	if game.Name.String() != gameName {
		t.Fatal("Bad name [got]", game.Name, "[want]", gameName)
	}
	if game.Description.String() != desc {
		t.Fatal("Bad description [got]", game.Description, "[want]", desc)
	}

	// Try to create duplicate
	_, err = service.Create(&dto.CreateGameDTO{
		Name: gameName,
	})
	if err == nil {
		t.Fatal("Error, you can't create duplicate game")
	}
	if !errors.Is(err, er.GameExist) {
		t.Fatal(err)
	}

	// Delete game
	err = service.Delete(game.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func testDelete(t *testing.T) {
	service := NewService()
	gameName := "one"
	gameID := utils.NameToID(gameName)

	// Try to remove non-existing game
	err := service.Delete(gameID)
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = service.Create(&dto.CreateGameDTO{
		Name: gameName,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete game
	err = service.Delete(gameID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete game twice
	err = service.Delete(gameID)
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}
}
func testUpdate(t *testing.T) {
	service := NewService()
	gameName := []string{"one", "two"}
	desc := []string{"first description", "second description"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Try to update non-existing game
	_, err := service.Update(gameID[0], &dto.UpdateGameDTO{})
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	game, err := service.Create(&dto.CreateGameDTO{
		Name:        gameName[0],
		Description: desc[0],
	})
	if err != nil {
		t.Fatal(err)
	}
	if game.Name.String() != gameName[0] {
		t.Fatal("Bad name [got]", game.Name, "[want]", gameName[0])
	}
	if game.Description.String() != desc[0] {
		t.Fatal("Bad description [got]", game.Description, "[want]", desc[0])
	}

	// Update game
	game, err = service.Update(gameID[0], &dto.UpdateGameDTO{
		Name:        gameName[1],
		Description: desc[1],
	})
	if err != nil {
		t.Fatal(err)
	}
	if game.Name.String() != gameName[1] {
		t.Fatal("Bad name [got]", game.Name, "[want]", gameName[1])
	}
	if game.Description.String() != desc[1] {
		t.Fatal("Bad description [got]", game.Description, "[want]", desc[1])
	}

	// Delete game
	err = service.Delete(gameID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Try to update non-existing game
	_, err = service.Update(gameID[1], &dto.UpdateGameDTO{})
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}
}
func testList(t *testing.T) {
	service := NewService()
	gameName := []string{"B game", "A game"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Empty list
	items, err := service.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}

	// Create first game
	_, err = service.Create(&dto.CreateGameDTO{
		Name: gameName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// One game
	items, err = service.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatal("List should be with 1 element")
	}

	// Create second game
	_, err = service.Create(&dto.CreateGameDTO{
		Name: gameName[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sort by name
	items, err = service.List("name")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != gameName[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[1])
	}
	if items[1].Name.String() != gameName[0] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[0])
	}

	// Sort by name_desc
	items, err = service.List("name_desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != gameName[0] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[0])
	}
	if items[1].Name.String() != gameName[1] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[1])
	}

	// Sort by created date
	items, err = service.List("created")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != gameName[0] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[0])
	}
	if items[1].Name.String() != gameName[1] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[1])
	}

	// Sort by created_desc
	items, err = service.List("created_desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != gameName[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[1])
	}
	if items[1].Name.String() != gameName[0] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[0])
	}

	// Delete first game
	err = service.Delete(gameID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Delete second game
	err = service.Delete(gameID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Empty list
	items, err = service.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}
}
func testItem(t *testing.T) {
	service := NewService()
	gameName := []string{"one", "two"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Try to get non-existing game
	_, err := service.Item(gameID[0])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = service.Create(&dto.CreateGameDTO{
		Name: gameName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid game
	_, err = service.Item(gameID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid game
	_, err = service.Item(gameID[1])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Rename game
	_, err = service.Update(gameID[0], &dto.UpdateGameDTO{Name: gameName[1]})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid game
	_, err = service.Item(gameID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid game
	_, err = service.Item(gameID[0])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Delete game
	err = service.Delete(gameID[1])
	if err != nil {
		t.Fatal(err)
	}
}
func testImage(t *testing.T) {
	service := NewService()
	gameName := "one"
	gameID := utils.NameToID(gameName)
	pngImage := "https://github.com/fluidicon.png"
	jpegImage := "https://avatars.githubusercontent.com/apple"

	// Check no game
	_, _, err := service.GetImage(gameID)
	if err == nil {
		t.Fatal("Error, game not exists")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = service.Create(&dto.CreateGameDTO{
		Name:  gameName,
		Image: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := service.GetImage(gameID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update game
	_, err = service.Update(gameID, &dto.UpdateGameDTO{
		Image: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = service.GetImage(gameID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update game
	_, err = service.Update(gameID, &dto.UpdateGameDTO{
		Image: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = service.GetImage(gameID)
	if err == nil {
		t.Fatal("Error, game don't have image")
	}
	if !errors.Is(err, er.GameImageNotExists) {
		t.Fatal(err)
	}

	// Delete game
	err = service.Delete(gameID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGame(t *testing.T) {
	t.Parallel()

	// Set path for the game test artifacts
	dataPath := os.Getenv("TEST_DATA_PATH")
	if dataPath == "" {
		t.Fatal("TEST_DATA_PATH must be set")
	}
	config.GetConfig().SetDataPath(filepath.Join(dataPath, "game_test"))

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
func fuzzList(t *testing.T, service *GameService, waitItems int) error {
	items, err := service.List("")
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
func fuzzItem(t *testing.T, service *GameService, gameID, name, desc string) error {
	game, err := service.Item(gameID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	if game.Name.String() != name {
		{
			data, _ := json.MarshalIndent(game, "", "	")
			t.Log(string(data))
		}
		return fmt.Errorf("name: [wait] %s [got] %s", name, game.Name)
	}
	if game.Description.String() != desc {
		{
			data, _ := json.MarshalIndent(game, "", "	")
			t.Log("item:", string(data))
		}
		return fmt.Errorf("description: [wait] %q [got] %q", desc, game.Description)
	}
	return nil
}
func fuzzCreate(t *testing.T, service *GameService, name, desc string) (*entity.GameInfo, error) {
	game, err := service.Create(&dto.CreateGameDTO{
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
func fuzzUpdate(t *testing.T, service *GameService, gameID, name, desc string) (*entity.GameInfo, error) {
	game, err := service.Update(gameID, &dto.UpdateGameDTO{
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
func fuzzDelete(t *testing.T, service *GameService, gameID string) error {
	err := service.Delete(gameID)
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
	// Set path for the game test artifacts
	dataPath := os.Getenv("TEST_DATA_PATH")
	if dataPath == "" {
		f.Fatal("TEST_DATA_PATH must be set")
	}
	config.GetConfig().SetDataPath(filepath.Join(dataPath, "game_fuzz_"+uuid.New().String()))

	service := NewService()

	msync := sync.Mutex{}
	f.Fuzz(func(t *testing.T, name1, desc1, name2, desc2 string) {
		if utils.NameToID(name1) == "" || utils.NameToID(name2) == "" {
			// skip
			return
		}

		// Only one test at once
		msync.Lock()
		defer msync.Unlock()

		// Empty list
		err := fuzzList(t, service, 0)
		if err != nil {
			t.Fatal(err)
		}

		// Create game
		game1, err := fuzzCreate(t, service, name1, desc1)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// List with game
		err = fuzzList(t, service, 1)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = fuzzItem(t, service, game1.ID, name1, desc1)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Update game
		game2, err := fuzzUpdate(t, service, utils.NameToID(name1), name2, desc2)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// List with game
		err = fuzzList(t, service, 1)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = fuzzItem(t, service, game2.ID, name2, desc2)
		if err != nil {
			fuzzCleanup(dataPath) // Cleanup - just in case
			t.Fatal(err)
		}

		// Delete game
		err = fuzzDelete(t, service, utils.NameToID(name2))
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
