package games

import (
	"io/ioutil"
	"log"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
)

type ListOfGamesResponse struct {
	Games []*GameInfo `json:"games"`
}

func ListOfGames() (result *ListOfGamesResponse, e *errors.Error) {
	result = &ListOfGamesResponse{
		Games: make([]*GameInfo, 0),
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
		var item *GameInfo

		// Check if game info exist
		var exist bool
		exist, e = GameIsInfoExist(file.Name())
		if e != nil {
			log.Println("No info:", file.Name())
			continue
		}
		if !exist {
			e = errors.GameInfoNotExists
			return
		}

		// Get info
		item, e = GameGetInfo(file.Name())
		if e != nil {
			log.Println("Bad info:", file.Name())
			continue
		}

		// Append game info to list
		result.Games = append(result.Games, item)
	}
	return
}
