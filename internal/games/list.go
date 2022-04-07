package games

import (
	"io/ioutil"
	"log"

	"tts_deck_build/internal/config"
)

type ListOfGamesResponse struct {
	Games []string `json:"games"`
}

func ListOfGames() (result *ListOfGamesResponse) {
	result = &ListOfGamesResponse{
		Games: make([]string, 0),
	}

	files, err := ioutil.ReadDir(config.GetConfig().Games())
	if err != nil {
		log.Println(err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			result.Games = append(result.Games, file.Name())
		}
	}
	return
}
