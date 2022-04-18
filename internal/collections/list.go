package collections

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/games"
)

type ListOfCollectionsResponse struct {
	Collections []*CollectionInfo `json:"collections"`
}

func ListOfCollections(gameName string) (result *ListOfCollectionsResponse, e *errors.Error) {
	e = games.FullGameCheck(gameName)
	if e != nil {
		return
	}
	result = &ListOfCollectionsResponse{
		Collections: make([]*CollectionInfo, 0),
	}

	gamePath := filepath.Join(config.GetConfig().Games(), gameName)
	files, err := ioutil.ReadDir(gamePath)
	if err != nil {
		e = errors.InternalError.AddMessage(err.Error())
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		var item *CollectionInfo

		// Check if collection info exist
		var exist bool
		exist, e = CollectionIsInfoExist(gameName, file.Name())
		if e != nil {
			log.Println("No info:", file.Name())
			continue
		}
		if !exist {
			e = errors.CollectionInfoNotExists
			return
		}

		// Get info
		item, e = CollectionGetInfo(gameName, file.Name())
		if e != nil {
			log.Println("Bad info:", file.Name())
			continue
		}

		// Append collection info to list
		result.Collections = append(result.Collections, item)
	}
	return
}
