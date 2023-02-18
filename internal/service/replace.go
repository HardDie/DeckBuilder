package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/HardDie/DeckBuilder/internal/tts_entity"
)

type IReplaceService interface {
	Prepare(data []byte) ([]Couple, error)
	Replace(data, mapping []byte) (*tts_entity.RootObjects, error)
}
type ReplaceService struct {
}

func NewReplaceService() *ReplaceService {
	return &ReplaceService{}
}

type Request struct {
	ObjectStates []struct {
		ContainedObjects []struct {
			CustomDeck map[string]struct {
				FaceURL string `json:"FaceURL"`
				BackURL string `json:"BackURL"`
			} `json:"CustomDeck"`
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
		return nil, errors.New("should be single root object")
	}

	var res []Couple
	uniq := make(map[string]string)
	for _, item := range req.ObjectStates[0].ContainedObjects {
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
			return fmt.Errorf("can't find mapping for %q back url", val.BackURL)
		}
		val.BackURL = newUrl

		// Replace front image
		newUrl, ok = mm[val.FaceURL]
		if !ok {
			return fmt.Errorf("can't find mapping for %q face url", val.FaceURL)
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
		return nil, fmt.Errorf("error parsing mapping file: %s", err.Error())
	}

	var root tts_entity.RootObjects
	err = json.Unmarshal(data, &root)
	if err != nil {
		return nil, fmt.Errorf("error parsing data file: %s", err.Error())
	}

	// Convert items into map
	mm := make(map[string]string)
	for _, val := range m.Data {
		mm[val.Key] = val.Value
	}

	if len(root.ObjectStates) != 1 {
		return nil, errors.New("should be single root object")
	}

	var newContained []any
	for _, item := range root.ObjectStates[0].ContainedObjects {
		tmp, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unknown object type: %T", item)
		}

		name, ok := tmp["Name"]
		if !ok {
			return nil, errors.New("object don't have Name field")
		}

		tmpJson, err := json.Marshal(tmp)
		if err != nil {
			return nil, fmt.Errorf("error marshaling %w", err)
		}

		switch name {
		case "Deck":
			var deck tts_entity.DeckObject
			err = json.Unmarshal(tmpJson, &deck)
			if err != nil {
				return nil, fmt.Errorf("error deck parsing %w", err)
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

			newContained = append(newContained, deck)
		case "Card":
			var card tts_entity.Card
			err = json.Unmarshal(tmpJson, &card)
			if err != nil {
				return nil, fmt.Errorf("error card parsing %w", err)
			}

			// Replace for custom deck
			err = replaceCustomDeck(card.CustomDeck, mm)
			if err != nil {
				return nil, err
			}

			newContained = append(newContained, card)
		default:
			return nil, fmt.Errorf("unknown object: %q", name)
		}
	}

	root.ObjectStates[0].ContainedObjects = newContained
	return &root, nil
}
