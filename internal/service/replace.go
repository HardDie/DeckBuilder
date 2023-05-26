package service

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/tts_entity"
)

type IReplaceService interface {
	Prepare(data []byte) ([]Couple, error)
	Replace(data, mapping []byte) (*tts_entity.RootObjects, error)
}
type ReplaceService struct {
	ttsService ITTSService
}

func NewReplaceService(ttsService ITTSService) *ReplaceService {
	return &ReplaceService{
		ttsService: ttsService,
	}
}

type Request struct {
	ObjectStates []struct {
		ContainedObjects []struct {
			ContainedObjects []struct {
				CustomDeck map[string]struct {
					FaceURL string `json:"FaceURL"`
					BackURL string `json:"BackURL"`
				} `json:"CustomDeck"`
			} `json:"containedObjects"`
		} `json:"ContainedObjects"`
	} `json:"ObjectStates"`
}
type Couple struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *ReplaceService) Prepare(data []byte) ([]Couple, error) {
	req := Request{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	if len(req.ObjectStates) != 1 {
		return nil, errors.ErrorInvalidDeckDescription
	}

	var res []Couple
	uniq := make(map[string]string)
	for _, collectionBag := range req.ObjectStates[0].ContainedObjects {
		for _, item := range collectionBag.ContainedObjects {
			for _, val := range item.CustomDeck {
				if _, ok := uniq[val.BackURL]; !ok {
					res = append(res, Couple{Key: val.BackURL})
					uniq[val.BackURL] = ""
				}
				if _, ok := uniq[val.FaceURL]; !ok {
					res = append(res, Couple{Key: val.FaceURL})
					uniq[val.FaceURL] = ""
				}
			}
		}
	}
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Key < res[j].Key
	})
	return res, nil
}

type Mapping struct {
	Data []Couple `json:"data"`
}

func replaceCustomDeck(customDeck map[int]tts_entity.DeckDescription, mm map[string]string) error {
	for key, val := range customDeck {
		// Replace back image
		newUrl, ok := mm[val.BackURL]
		if !ok {
			return errors.ErrorInvalidDeckDescription.AddMessage(fmt.Sprintf("can't find mapping for %q back url", val.BackURL))
		}
		val.BackURL = newUrl

		// Replace front image
		newUrl, ok = mm[val.FaceURL]
		if !ok {
			return errors.ErrorInvalidDeckDescription.AddMessage(fmt.Sprintf("can't find mapping for %q face url", val.FaceURL))
		}
		val.FaceURL = newUrl

		customDeck[key] = val
	}
	return nil
}

func (s *ReplaceService) Replace(data, mapping []byte) (*tts_entity.RootObjects, error) {
	var m Mapping
	err := json.Unmarshal(mapping, &m)
	if err != nil {
		logger.Info.Printf("error parsing mapping file: %s", err.Error())
		return nil, errors.ErrorInvalidMapping
	}

	var root tts_entity.RootObjects
	err = json.Unmarshal(data, &root)
	if err != nil {
		logger.Info.Printf("error parsing data file: %s", err.Error())
		return nil, errors.ErrorInvalidDeckDescription
	}

	// Convert items into map
	mm := make(map[string]string)
	for _, val := range m.Data {
		mm[val.Key] = val.Value
	}

	if len(root.ObjectStates) != 1 {
		logger.Info.Println("should be single root object")
		return nil, errors.ErrorInvalidDeckDescription
	}

	var newContained []any
	for _, collectionBagTemp := range root.ObjectStates[0].ContainedObjects {
		tmp, ok := collectionBagTemp.(map[string]any)
		if !ok {
			logger.Info.Printf("unknown object type: %T", collectionBagTemp)
			return nil, errors.ErrorInvalidDeckDescription
		}

		name, ok := tmp["Name"]
		if !ok {
			logger.Info.Println("object don't have Name field")
			return nil, errors.ErrorInvalidDeckDescription
		}
		if name != "Bag" {
			logger.Info.Println("error collection bag parsing %w", err)
			return nil, errors.ErrorInvalidDeckDescription
		}

		tmpJson, err := json.Marshal(tmp)
		if err != nil {
			logger.Info.Printf("error marshaling %w", err)
			return nil, errors.ErrorInvalidDeckDescription
		}

		var collectionBag tts_entity.Bag
		err = json.Unmarshal(tmpJson, &collectionBag)
		if err != nil {
			logger.Info.Println("error collection bag parsing %w", err)
			return nil, errors.ErrorInvalidDeckDescription
		}

		var collectionBagContaind []any
		for _, item := range collectionBag.ContainedObjects {
			tmp, ok := item.(map[string]any)
			if !ok {
				logger.Info.Printf("unknown object type: %T", item)
				return nil, errors.ErrorInvalidDeckDescription
			}

			name, ok := tmp["Name"]
			if !ok {
				logger.Info.Println("object don't have Name field")
				return nil, errors.ErrorInvalidDeckDescription
			}

			tmpJson, err := json.Marshal(tmp)
			if err != nil {
				logger.Info.Printf("error marshaling %w", err)
				return nil, errors.ErrorInvalidDeckDescription
			}

			switch name {
			case "Deck":
				var deck tts_entity.DeckObject
				err = json.Unmarshal(tmpJson, &deck)
				if err != nil {
					logger.Info.Printf("error deck parsing %w", err)
					return nil, errors.ErrorInvalidDeckDescription
				}

				// Replace for custom deck
				err = replaceCustomDeck(deck.CustomDeck, mm)
				if err != nil {
					return nil, err
				}

				// Replace for cards inside deck
				for i, card := range deck.ContainedObjects {
					err = replaceCustomDeck(card.CustomDeck, mm)
					if err != nil {
						return nil, err
					}
					deck.ContainedObjects[i] = card
				}

				collectionBagContaind = append(collectionBagContaind, deck)
			case "Card":
				var card tts_entity.Card
				err = json.Unmarshal(tmpJson, &card)
				if err != nil {
					logger.Info.Printf("error card parsing %w", err)
					return nil, errors.ErrorInvalidDeckDescription
				}

				// Replace for custom deck
				err = replaceCustomDeck(card.CustomDeck, mm)
				if err != nil {
					return nil, err
				}

				collectionBagContaind = append(collectionBagContaind, card)
			default:
				logger.Info.Printf("unknown object: %q", name)
				return nil, errors.ErrorInvalidDeckDescription
			}
		}
		if len(collectionBagContaind) > 0 {
			collectionBag.ContainedObjects = collectionBagContaind
			newContained = append(newContained, collectionBag)
		}
	}

	root.ObjectStates[0].ContainedObjects = newContained
	s.ttsService.SendToTTS(root.ObjectStates[0])

	return &root, nil
}
