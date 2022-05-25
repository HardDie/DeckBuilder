package games

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"tts_deck_build/internal/config"
	er "tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

func testCreate(t *testing.T) {
	service := NewService()
	gameName := "one"
	desc := "best game ever"

	// Create game
	game, err := service.Create(&CreateGameDTO{
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
	_, err = service.Create(&CreateGameDTO{
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
	gameId := utils.NameToID(gameName)

	// Try to remove non-existing game
	err := service.Delete(gameId)
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = service.Create(&CreateGameDTO{
		Name: gameName,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete game
	err = service.Delete(gameId)
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete game twice
	err = service.Delete(gameId)
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
	gameId := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Try to update non-existing game
	_, err := service.Update(gameId[0], &UpdateGameDTO{})
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	game, err := service.Create(&CreateGameDTO{
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
	game, err = service.Update(gameId[0], &UpdateGameDTO{
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
	err = service.Delete(gameId[1])
	if err != nil {
		t.Fatal(err)
	}

	// Try to update non-existing game
	_, err = service.Update(gameId[1], &UpdateGameDTO{})
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
	gameId := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Empty list
	items, err := service.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}

	// Create first game
	_, err = service.Create(&CreateGameDTO{
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
	_, err = service.Create(&CreateGameDTO{
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
	if items[0].Name != gameName[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[1])
	}
	if items[1].Name != gameName[0] {
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
	if items[0].Name != gameName[0] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[0])
	}
	if items[1].Name != gameName[1] {
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
	if items[0].Name != gameName[0] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[0])
	}
	if items[1].Name != gameName[1] {
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
	if items[0].Name != gameName[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[1])
	}
	if items[1].Name != gameName[0] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[0])
	}

	// Delete first game
	err = service.Delete(gameId[0])
	if err != nil {
		t.Fatal(err)
	}

	// Delete second game
	err = service.Delete(gameId[1])
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
	gameId := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Try to get non-existing game
	_, err := service.Item(gameId[0])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = service.Create(&CreateGameDTO{
		Name: gameName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid game
	_, err = service.Item(gameId[0])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid game
	_, err = service.Item(gameId[1])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Rename game
	_, err = service.Update(gameId[0], &UpdateGameDTO{Name: gameName[1]})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid game
	_, err = service.Item(gameId[1])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid game
	_, err = service.Item(gameId[0])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Delete game
	err = service.Delete(gameId[1])
	if err != nil {
		t.Fatal(err)
	}
}
func testImage(t *testing.T) {
	service := NewService()
	gameName := "one"
	gameId := utils.NameToID(gameName)
	pngImage := "https://github.com/fluidicon.png"
	jpegImage := "https://avatars.githubusercontent.com/apple"

	// Create game
	_, err := service.Create(&CreateGameDTO{
		Name:  gameName,
		Image: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := service.GetImage(gameId)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update game
	_, err = service.Update(gameId, &UpdateGameDTO{
		Image: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = service.GetImage(gameId)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Delete game
	err = service.Delete(gameId)
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
