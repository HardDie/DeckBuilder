package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/tts_entity"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

func readBag(path string) RootRaw {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error read file %q: %s", path, err.Error())
	}

	var result RootRaw
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("error unmarshal file %q: %s", path, err.Error())
	}

	return result
}

type RootRaw struct {
	ObjectStates []GameBagRaw `json:"ObjectStates"`
}

type GameBagRaw struct {
	Name             string               `json:"Name"`
	Transform        tts_entity.Transform `json:"Transform"`
	Nickname         string               `json:"Nickname"`
	Description      string               `json:"Description"`
	ContainedObjects []CollectionBagRaw   `json:"ContainedObjects"`
}

type CollectionBagRaw struct {
	Name             string               `json:"Name"`
	Transform        tts_entity.Transform `json:"Transform"`
	Nickname         string               `json:"Nickname"`
	Description      string               `json:"Description"`
	ContainedObjects []json.RawMessage    `json:"ContainedObjects"`
}

func RawToValid(bag RootRaw) Root {
	if len(bag.ObjectStates) != 1 {
		log.Fatal("root must have 1 object")
	}
	gameBagRaw := bag.ObjectStates[0]

	gameBag := GameBag{
		Name:        gameBagRaw.Name,
		Transform:   gameBagRaw.Transform,
		Nickname:    gameBagRaw.Nickname,
		Description: gameBagRaw.Description,
	}

	for _, colRaw := range gameBagRaw.ContainedObjects {
		col := CollectionBag{
			Name:        colRaw.Name,
			Transform:   colRaw.Transform,
			Nickname:    colRaw.Nickname,
			Description: colRaw.Description,
		}

		for _, itemRaw := range colRaw.ContainedObjects {
			name, err := getName(itemRaw)
			if err != nil {
				log.Fatal("error get name of raw item:", err.Error())
			}
			switch name {
			case "Deck":
				var deck tts_entity.DeckObject
				err = utils.ObjectJSONObject(itemRaw, &deck)
				if err != nil {
					log.Fatal("error parse deck:", err.Error())
				}
				col.ContainedObjects = append(col.ContainedObjects, deck)
			case "Card":
				var card tts_entity.Card
				err = utils.ObjectJSONObject(itemRaw, &card)
				if err != nil {
					log.Fatal("error parse card:", err.Error())
				}
				col.ContainedObjects = append(col.ContainedObjects, card)
			default:
				log.Fatal("unknown item type:", name)
			}
		}
		gameBag.ContainedObjects = append(gameBag.ContainedObjects, col)
	}

	return Root{
		ObjectStates: []GameBag{gameBag},
	}
}

func getName(obj json.RawMessage) (string, error) {
	var tmp map[string]any
	err := json.Unmarshal(obj, &tmp)
	if err != nil {
		return "", err
	}
	name, ok := tmp["Name"]
	if !ok {
		logger.Info.Println("object don't have Name field")
		return "", errors.ErrorInvalidDeckDescription
	}
	nameStr, ok := name.(string)
	if !ok {
		logger.Info.Println("Name field is not string")
		return "", errors.ErrorInvalidDeckDescription
	}
	return nameStr, nil
}
