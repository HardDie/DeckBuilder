package replace

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"

	"github.com/HardDie/DeckBuilder/internal/apperr"
	servicesTTS "github.com/HardDie/DeckBuilder/internal/services/tts"
	"github.com/HardDie/DeckBuilder/internal/tts_entity"
	"github.com/HardDie/DeckBuilder/internal/utils"
)

type replace struct {
	serviceTTS servicesTTS.TTS
}

func New(serviceTTS servicesTTS.TTS) Replace {
	return &replace{
		serviceTTS: serviceTTS,
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

func (s *replace) Prepare(data []byte) ([]Couple, error) {
	req := Request{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		slog.Info("render file is not JSON", "err", err)
		return nil, apperr.ErrBadRenderFile
	}

	if len(req.ObjectStates) != 1 {
		return nil, apperr.ErrBadRenderFile
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
			return apperr.Withf(apperr.ErrBadMappingFile, "the mapping file has no URL for %q", val.BackURL)
		}
		val.BackURL = newUrl

		// Replace front image
		newUrl, ok = mm[val.FaceURL]
		if !ok {
			return apperr.Withf(apperr.ErrBadMappingFile, "the mapping file has no URL for %q", val.FaceURL)
		}
		val.FaceURL = newUrl

		customDeck[key] = val
	}
	return nil
}

func (s *replace) Replace(data, mapping []byte) (*tts_entity.RootObjects, error) {
	var m Mapping
	err := json.Unmarshal(mapping, &m)
	if err != nil {
		slog.Info("mapping file is not JSON", "err", err)
		return nil, apperr.ErrBadMappingFile
	}

	var root tts_entity.RootObjects
	err = json.Unmarshal(data, &root)
	if err != nil {
		slog.Info("render file is not JSON", "err", err)
		return nil, apperr.ErrBadRenderFile
	}

	// Convert items into map
	mm := make(map[string]string)
	for _, val := range m.Data {
		mm[val.Key] = val.Value
	}

	if len(root.ObjectStates) != 1 {
		slog.Info("render file must have one root object", "objects", len(root.ObjectStates))
		return nil, apperr.ErrBadRenderFile
	}

	var newContained []any
	for _, collectionBagTemp := range root.ObjectStates[0].ContainedObjects {
		colName, err := getName(collectionBagTemp)
		if err != nil {
			slog.Info("render file: collection without a name", "type", fmt.Sprintf("%T", collectionBagTemp))
			return nil, err
		}
		if colName != "Bag" {
			slog.Info("render file: bad collection bag", "err", err)
			return nil, apperr.ErrBadRenderFile
		}

		var collectionBag tts_entity.Bag
		err = utils.ObjectJSONObject(collectionBagTemp, &collectionBag)
		if err != nil {
			slog.Info("render file: bad collection bag", "err", err)
			return nil, apperr.ErrBadRenderFile
		}

		var collectionBagContaind []any
		for _, item := range collectionBag.ContainedObjects {
			itemName, err := getName(item)
			if err != nil {
				slog.Info("render file: item without a name", "type", fmt.Sprintf("%T", item))
				return nil, err
			}

			switch itemName {
			case "Deck":
				var deck tts_entity.DeckObject
				err = utils.ObjectJSONObject(item, &deck)
				if err != nil {
					slog.Info("render file: bad deck", "err", err)
					return nil, apperr.ErrBadRenderFile
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
				err = utils.ObjectJSONObject(item, &card)
				if err != nil {
					slog.Info("render file: bad card", "err", err)
					return nil, apperr.ErrBadRenderFile
				}

				// Replace for custom deck
				err = replaceCustomDeck(card.CustomDeck, mm)
				if err != nil {
					return nil, err
				}

				collectionBagContaind = append(collectionBagContaind, card)
			default:
				slog.Info("render file: unknown object", "name", itemName)
				return nil, apperr.ErrBadRenderFile
			}
		}
		if len(collectionBagContaind) > 0 {
			collectionBag.ContainedObjects = collectionBagContaind
			newContained = append(newContained, collectionBag)
		}
	}

	root.ObjectStates[0].ContainedObjects = newContained
	s.serviceTTS.SendToTTS(root.ObjectStates[0])

	return &root, nil
}

func getName(obj any) (string, error) {
	tmp, ok := obj.(map[string]any)
	if !ok {
		slog.Info("render file: unknown object type", "type", fmt.Sprintf("%T", obj))
		return "", apperr.ErrBadRenderFile
	}
	name, ok := tmp["Name"]
	if !ok {
		slog.Info("render file: object has no Name field")
		return "", apperr.ErrBadRenderFile
	}
	nameStr, ok := name.(string)
	if !ok {
		slog.Info("render file: Name field is not a string")
		return "", apperr.ErrBadRenderFile
	}
	return nameStr, nil
}
