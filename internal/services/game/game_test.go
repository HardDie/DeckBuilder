package game

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/HardDie/fsentry"
	"github.com/stretchr/testify/assert"

	"github.com/HardDie/DeckBuilder/internal/config"
	dbCore "github.com/HardDie/DeckBuilder/internal/db/core"
	dbGame "github.com/HardDie/DeckBuilder/internal/db/game"
	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
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
	assert.NoError(t, err)
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
	g, err := tt.serviceGame.Create(CreateRequest{
		Name:        gameName,
		Description: desc,
	})
	assert.NoError(t, err)
	assert.Equal(t, gameName, g.Name)
	assert.Equal(t, desc, g.Description)

	// Try to create duplicate
	_, err = tt.serviceGame.Create(CreateRequest{
		Name: gameName,
	})
	assert.ErrorIs(t, err, er.GameExist)

	// Delete game
	err = tt.serviceGame.Delete(g.ID)
	assert.NoError(t, err)
}
func (tt *gameTest) testDelete(t *testing.T) {
	gameName := "delete_one"
	gameID := utils.NameToID(gameName)

	// Try to remove non-existing game
	err := tt.serviceGame.Delete(gameID)
	assert.ErrorIs(t, err, er.GameNotExists)

	// Create game
	_, err = tt.serviceGame.Create(CreateRequest{
		Name: gameName,
	})
	assert.NoError(t, err)

	// Delete game
	err = tt.serviceGame.Delete(gameID)
	assert.NoError(t, err)

	// Try to delete game twice
	err = tt.serviceGame.Delete(gameID)
	assert.ErrorIs(t, err, er.GameNotExists)
}
func (tt *gameTest) testUpdate(t *testing.T) {
	gameName := []string{"update_one", "update_two"}
	desc := []string{"first description", "second description"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Try to update non-existing game
	_, err := tt.serviceGame.Update(gameID[0], UpdateRequest{})
	assert.ErrorIs(t, err, er.GameNotExists)

	// Create game
	g, err := tt.serviceGame.Create(CreateRequest{
		Name:        gameName[0],
		Description: desc[0],
	})
	assert.NoError(t, err)
	assert.Equal(t, gameName[0], g.Name)
	assert.Equal(t, desc[0], g.Description)

	// Update game
	g, err = tt.serviceGame.Update(gameID[0], UpdateRequest{
		Name:        gameName[1],
		Description: desc[1],
	})
	assert.NoError(t, err)
	assert.Equal(t, gameName[1], g.Name)
	assert.Equal(t, desc[1], g.Description)

	// Delete game
	err = tt.serviceGame.Delete(gameID[1])
	assert.NoError(t, err)

	// Try to update non-existing game
	_, err = tt.serviceGame.Update(gameID[1], UpdateRequest{})
	assert.ErrorIs(t, err, er.GameNotExists)
}
func (tt *gameTest) testList(t *testing.T) {
	gameName := []string{"B game", "A game"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Empty list
	items, err := tt.serviceGame.List("", "")
	assert.NoError(t, err)
	assert.Len(t, items, 0)

	// Create first game
	_, err = tt.serviceGame.Create(CreateRequest{
		Name: gameName[0],
	})
	assert.NoError(t, err)

	// One game
	items, err = tt.serviceGame.List("", "")
	assert.NoError(t, err)
	assert.Len(t, items, 1)

	// Create second game
	_, err = tt.serviceGame.Create(CreateRequest{
		Name: gameName[1],
	})
	assert.NoError(t, err)

	// Sort by name
	items, err = tt.serviceGame.List("name", "")
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, gameName[1], items[0].Name)
	assert.Equal(t, gameName[0], items[1].Name)

	// Sort by name_desc
	items, err = tt.serviceGame.List("name_desc", "")
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, gameName[0], items[0].Name)
	assert.Equal(t, gameName[1], items[1].Name)

	// Sort by created date
	items, err = tt.serviceGame.List("created", "")
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, gameName[0], items[0].Name)
	assert.Equal(t, gameName[1], items[1].Name)

	// Sort by created_desc
	items, err = tt.serviceGame.List("created_desc", "")
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, gameName[1], items[0].Name)
	assert.Equal(t, gameName[0], items[1].Name)

	// Delete first game
	err = tt.serviceGame.Delete(gameID[0])
	assert.NoError(t, err)

	// Delete second game
	err = tt.serviceGame.Delete(gameID[1])
	assert.NoError(t, err)

	// Empty list
	items, err = tt.serviceGame.List("", "")
	assert.NoError(t, err)
	assert.Len(t, items, 0)
}
func (tt *gameTest) testItem(t *testing.T) {
	gameName := []string{"item_one", "item_two"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Try to get non-existing game
	_, err := tt.serviceGame.Item(gameID[0])
	assert.ErrorIs(t, err, er.GameNotExists)

	// Create game
	_, err = tt.serviceGame.Create(CreateRequest{
		Name: gameName[0],
	})
	assert.NoError(t, err)

	// Get valid game
	_, err = tt.serviceGame.Item(gameID[0])
	assert.NoError(t, err)

	// Get invalid game
	_, err = tt.serviceGame.Item(gameID[1])
	assert.ErrorIs(t, err, er.GameNotExists)

	// Rename game
	_, err = tt.serviceGame.Update(gameID[0], UpdateRequest{
		Name: gameName[1],
	})
	assert.NoError(t, err)

	// Get valid game
	_, err = tt.serviceGame.Item(gameID[1])
	assert.NoError(t, err)

	// Get invalid game
	_, err = tt.serviceGame.Item(gameID[0])
	assert.ErrorIs(t, err, er.GameNotExists)

	// Delete game
	err = tt.serviceGame.Delete(gameID[1])
	assert.NoError(t, err)
}
func (tt *gameTest) testDuplicate(t *testing.T) {
	gameName := []string{"duplicate_one", "duplicate_two"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Create games
	_, err := tt.serviceGame.Create(CreateRequest{
		Name: gameName[0],
	})
	assert.NoError(t, err)
	_, err = tt.serviceGame.Create(CreateRequest{
		Name: gameName[1],
	})
	assert.NoError(t, err)

	// Try to duplicate not exist game
	_, err = tt.serviceGame.Duplicate("not_exist_game", DuplicateRequest{
		Name: "new_game",
	})
	assert.ErrorIs(t, err, er.GameNotExists)

	// Try to duplicate to exist game
	_, err = tt.serviceGame.Duplicate(gameID[0], DuplicateRequest{
		Name: gameID[1],
	})
	assert.ErrorIs(t, err, er.GameExist)

	_, err = tt.serviceGame.Duplicate(gameID[0], DuplicateRequest{
		Name: "good_duplicate",
	})
	assert.NoError(t, err)

	err = tt.serviceGame.Delete(gameID[0])
	assert.NoError(t, err)
	err = tt.serviceGame.Delete(gameID[1])
	assert.NoError(t, err)
	err = tt.serviceGame.Delete("good_duplicate")
	assert.NoError(t, err)
}
func (tt *gameTest) testImage(t *testing.T) {
	gameName := "image_one"
	gameID := utils.NameToID(gameName)
	pngImage := "https://github.com/fluidicon.png"
	jpegImage := "https://avatars.githubusercontent.com/apple"

	// Check no game
	_, _, err := tt.serviceGame.GetImage(gameID)
	assert.ErrorIs(t, err, er.GameNotExists)

	// Create game
	_, err = tt.serviceGame.Create(CreateRequest{
		Name:  gameName,
		Image: pngImage,
	})
	assert.NoError(t, err)

	// Check image type
	_, imgType, err := tt.serviceGame.GetImage(gameID)
	assert.NoError(t, err)
	assert.Equal(t, imgType, "png")

	// Update game
	_, err = tt.serviceGame.Update(gameID, UpdateRequest{
		Name:  gameName,
		Image: jpegImage,
	})
	assert.NoError(t, err)

	// Check image type
	_, imgType, err = tt.serviceGame.GetImage(gameID)
	assert.NoError(t, err)
	assert.Equal(t, imgType, "jpeg")

	// Update game
	_, err = tt.serviceGame.Update(gameID, UpdateRequest{
		Name:  gameName,
		Image: "",
	})
	assert.NoError(t, err)

	// Check no image
	_, _, err = tt.serviceGame.GetImage(gameID)
	assert.ErrorIs(t, err, er.GameImageNotExists)

	// Delete game
	err = tt.serviceGame.Delete(gameID)
	assert.NoError(t, err)
}
func (tt *gameTest) testImageBin(t *testing.T) {
	gameName := "image_bin_one"
	gameID := utils.NameToID(gameName)

	pageImage := images.CreateImage(100, 100)
	pngImage, err := images.ImageToPng(pageImage)
	assert.NoError(t, err)
	jpegImage, err := images.ImageToJpeg(pageImage)
	assert.NoError(t, err)
	gifImage, err := images.ImageToGif(pageImage)
	assert.NoError(t, err)

	// Check no game
	_, _, err = tt.serviceGame.GetImage(gameID)
	assert.ErrorIs(t, err, er.GameNotExists)

	// Create game
	_, err = tt.serviceGame.Create(CreateRequest{
		Name:      gameName,
		ImageFile: pngImage,
	})
	assert.NoError(t, err)

	// Check image type
	_, imgType, err := tt.serviceGame.GetImage(gameID)
	assert.NoError(t, err)
	assert.Equal(t, imgType, "png")

	// Update game
	_, err = tt.serviceGame.Update(gameID, UpdateRequest{
		Name:      gameName,
		ImageFile: jpegImage,
	})
	assert.NoError(t, err)

	// Check image type
	_, imgType, err = tt.serviceGame.GetImage(gameID)
	assert.NoError(t, err)
	assert.Equal(t, imgType, "jpeg")

	// Update game
	_, err = tt.serviceGame.Update(gameID, UpdateRequest{
		Name:      gameName,
		ImageFile: gifImage,
	})
	assert.NoError(t, err)

	// Check image type
	_, imgType, err = tt.serviceGame.GetImage(gameID)
	assert.NoError(t, err)
	assert.Equal(t, imgType, "gif")

	// Update game
	_, err = tt.serviceGame.Update(gameID, UpdateRequest{
		Name: gameName,
	})
	assert.NoError(t, err)

	// Check image type
	_, imgType, err = tt.serviceGame.GetImage(gameID)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.Equal(t, imgType, "gif")

	// Update game
	_, err = tt.serviceGame.Update(gameID, UpdateRequest{
		Name:  gameName,
		Image: "empty",
	})
	assert.NoError(t, err)

	// Check no image
	_, _, err = tt.serviceGame.GetImage(gameID)
	assert.ErrorIs(t, err, er.GameImageNotExists)

	// Delete game
	err = tt.serviceGame.Delete(gameID)
	assert.NoError(t, err)
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
	items, err := tt.serviceGame.List("", "")
	assert.NoError(t, err)
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
func (tt *gameTest) fuzzCreate(t *testing.T, name, desc string) (*entitiesGame.Game, error) {
	g, err := tt.serviceGame.Create(CreateRequest{
		Name:        name,
		Description: desc,
	})
	assert.NoError(t, err)
	{
		data, _ := json.MarshalIndent(g, "", "	")
		t.Log("create:", string(data))
	}
	return g, nil
}
func (tt *gameTest) fuzzUpdate(t *testing.T, gameID, name, desc string) (*entitiesGame.Game, error) {
	g, err := tt.serviceGame.Update(gameID, UpdateRequest{
		Name:        name,
		Description: desc,
	})
	assert.NoError(t, err)
	{
		data, _ := json.MarshalIndent(g, "", "	")
		t.Log("update:", string(data))
	}
	return g, nil
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
