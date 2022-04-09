package games

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/utils"
)

type ListOfGamesResponse struct {
	Games []GameInfo `json:"games"`
}

func parseFile(infoFile string) (item GameInfo, err error) {
	// Open file
	f, err := os.Open(infoFile)
	if err != nil {
		utils.IfErrorLog(err)
		return
	}
	defer func() { utils.IfErrorLog(f.Close()) }()

	// Decode json
	err = json.NewDecoder(f).Decode(&item)
	if err != nil {
		utils.IfErrorLog(err)
		return
	}
	return
}

func ListOfGames() (result *ListOfGamesResponse, e *errors.Error) {
	result = &ListOfGamesResponse{
		Games: make([]GameInfo, 0),
	}

	files, err := ioutil.ReadDir(config.GetConfig().Games())
	if err != nil {
		e = errors.InternalError.AddMessage(err.Error())
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		var item GameInfo
		infoFile := filepath.Join(config.GetConfig().Games(), file.Name(), GameInfoFilename)

		// Check if game info exist
		_, err = os.Stat(infoFile)
		if os.IsNotExist(err) {
			log.Println("No info", infoFile)
			continue
		}

		// Parse file
		item, err = parseFile(infoFile)
		if err != nil {
			continue
		}

		// Append game to list
		result.Games = append(result.Games, item)
	}
	return
}
