package games

import (
	"os"
	"path/filepath"
	"testing"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/utils"
)

func TestGame(t *testing.T) {
	t.Parallel()

	dataPath := os.Getenv("TEST_DATA_PATH")
	if dataPath == "" {
		t.Fatal("TEST_DATA_PATH must be set")
	}

	// Set path for the game test artifacts
	config.GetConfig().SetDataPath(filepath.Join(dataPath, "game_test"))

	service := NewService()

	gameName := []string{"B game", "A game", "C game"}
	gameId := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1]), utils.NameToID(gameName[2])}

	t.Run("[list] no game", func(t *testing.T) {
		items, err := NewService().List("")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 0 {
			t.Fatal("List should be empty")
		}
	})

	// Prepare: create first game
	_, err := service.Create(&CreateGameDTO{
		Name:        gameName[0],
		Description: "First game",
	})
	if err != nil {
		t.Fatal()
	}

	t.Run("[list] single game", func(t *testing.T) {
		items, err := NewService().List("")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 1 {
			t.Fatal("List should with 1 value")
		}
	})

	t.Run("[item] single game", func(t *testing.T) {
		item, err := NewService().Item(gameId[0])
		if err != nil {
			t.Fatal(err)
		}
		if item.Name != gameName[0] {
			t.Fatal("Got wrong value: [got]", item.Name, "[want]", gameName[0])
		}
	})

	t.Run("[item] single game bad id", func(t *testing.T) {
		_, err = NewService().Item(gameId[1])
		if err == nil {
			t.Fatal("Should be error, bad id")
		}
	})

	// Prepare: create second game
	_, err = service.Create(&CreateGameDTO{
		Name:        gameName[1],
		Description: "Second game",
	})
	if err != nil {
		t.Fatal()
	}

	t.Run("[list] two games name sort", func(t *testing.T) {
		items, err := NewService().List("name")
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
	})

	t.Run("[list] two games name_desc sort", func(t *testing.T) {
		items, err := NewService().List("name_desc")
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
	})

	t.Run("[list] two games created sort", func(t *testing.T) {
		items, err := NewService().List("created")
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
	})

	t.Run("[list] two games created_desc sort", func(t *testing.T) {
		items, err := NewService().List("created_desc")
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
	})

	t.Run("[item] two games first", func(t *testing.T) {
		item, err := NewService().Item(gameId[0])
		if err != nil {
			t.Fatal(err)
		}
		if item.Name != gameName[0] {
			t.Fatal("Got wrong value: [got]", item.Name, "[want]", gameName[0])
		}
	})

	t.Run("[item] two games second", func(t *testing.T) {
		item, err := NewService().Item(gameId[1])
		if err != nil {
			t.Fatal(err)
		}
		if item.Name != gameName[1] {
			t.Fatal("Got wrong value: [got]", item.Name, "[want]", gameName[1])
		}
	})

	// Prepare: rename first game
	_, err = service.Update(gameId[0], &UpdateGameDTO{
		Name: gameName[2],
	})

	t.Run("[list] two games name sort", func(t *testing.T) {
		items, err := NewService().List("name")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 2 {
			t.Fatal("List should with 2 value")
		}

		if items[0].Name != gameName[1] {
			t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[1])
		}
		if items[1].Name != gameName[2] {
			t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[2])
		}
	})

	t.Run("[list] two games name_desc sort", func(t *testing.T) {
		items, err := NewService().List("name_desc")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 2 {
			t.Fatal("List should with 2 value")
		}

		if items[0].Name != gameName[2] {
			t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[2])
		}
		if items[1].Name != gameName[1] {
			t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[1])
		}
	})

	t.Run("[list] two games created sort", func(t *testing.T) {
		items, err := NewService().List("created")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 2 {
			t.Fatal("List should with 2 value")
		}

		if items[0].Name != gameName[2] {
			t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[2])
		}
		if items[1].Name != gameName[1] {
			t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[1])
		}
	})

	t.Run("[list] two games created_desc sort", func(t *testing.T) {
		items, err := NewService().List("created_desc")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 2 {
			t.Fatal("List should with 2 value")
		}

		if items[0].Name != gameName[1] {
			t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[1])
		}
		if items[1].Name != gameName[2] {
			t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[2])
		}
	})

	t.Run("[item] two games first", func(t *testing.T) {
		item, err := NewService().Item(gameId[2])
		if err != nil {
			t.Fatal(err)
		}
		if item.Name != gameName[2] {
			t.Fatal("Got wrong value: [got]", item.Name, "[want]", gameName[2])
		}
	})

	t.Run("[item] two games second", func(t *testing.T) {
		item, err := NewService().Item(gameId[1])
		if err != nil {
			t.Fatal(err)
		}
		if item.Name != gameName[1] {
			t.Fatal("Got wrong value: [got]", item.Name, "[want]", gameName[1])
		}
	})

	t.Run("[item] two games bad id after update", func(t *testing.T) {
		_, err = NewService().Item(gameId[0])
		if err == nil {
			t.Fatal("Should be error, bad id")
		}
	})

	// Prepare: delete first game
	err = service.Delete(gameId[2])
	if err != nil {
		t.Fatal()
	}

	t.Run("[list] single game after delete", func(t *testing.T) {
		items, err := NewService().List("")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 1 {
			t.Fatal("List should with 1 value")
		}
	})

	t.Run("[item] single game after delete", func(t *testing.T) {
		item, err := NewService().Item(gameId[1])
		if err != nil {
			t.Fatal(err)
		}
		if item.Name != gameName[1] {
			t.Fatal("Got wrong value: [got]", item.Name, "[want]", gameName[1])
		}
	})

	t.Run("[item] single game bad id after delete", func(t *testing.T) {
		_, err = NewService().Item(gameId[0])
		if err == nil {
			t.Fatal("Should be error, bad id")
		}
	})

	// Prepare: delete second game
	err = service.Delete(gameId[1])
	if err != nil {
		t.Fatal()
	}

	t.Run("[list] no game after delete", func(t *testing.T) {
		items, err := NewService().List("")
		if err != nil {
			t.Fatal(err)
		}
		if len(items) != 0 {
			t.Fatal("List should be empty")
		}
	})

	t.Run("[item] no games bad first id after delete", func(t *testing.T) {
		_, err = NewService().Item(gameId[2])
		if err == nil {
			t.Fatal("Should be error, bad id")
		}
	})

	t.Run("[item] no games bad second id after delete", func(t *testing.T) {
		_, err = NewService().Item(gameId[1])
		if err == nil {
			t.Fatal("Should be error, bad id")
		}
	})
}
