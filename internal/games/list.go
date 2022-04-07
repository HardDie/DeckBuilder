package games

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
)

type ListOfGamesResponse struct {
	Games []GameInfo `json:"games"`
}

func ListOfGames() (result *ListOfGamesResponse) {
	result = &ListOfGamesResponse{
		Games: make([]GameInfo, 0),
	}

	files, err := ioutil.ReadDir(config.GetConfig().Games())
	if err != nil {
		log.Println(err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			infoFile := filepath.Join(config.GetConfig().Games(), file.Name(), GameInfoFilename)
			_, err = os.Stat(infoFile)
			if os.IsNotExist(err) {
				log.Println("No info", infoFile)
				continue
			}
			data, err := ioutil.ReadFile(infoFile)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			item := GameInfo{}
			err = json.Unmarshal(data, &item)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			result.Games = append(result.Games, item)
		}
	}
	return
}
