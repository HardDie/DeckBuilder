package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/HardDie/DeckBuilder/internal/tts_entity"
)

type Root struct {
	ObjectStates []GameBag `json:"ObjectStates"`
}

type GameBag struct {
	Name             string               `json:"Name"`
	Transform        tts_entity.Transform `json:"Transform"`
	Nickname         string               `json:"Nickname"`
	Description      string               `json:"Description"`
	ContainedObjects []CollectionBag      `json:"ContainedObjects"`
}

type CollectionBag struct {
	Name             string                  `json:"Name"`
	Transform        tts_entity.Transform    `json:"Transform"`
	Nickname         string                  `json:"Nickname"`
	Description      string                  `json:"Description"`
	ContainedObjects []tts_entity.DeckObject `json:"ContainedObjects"`
}

func readBag(path string) Root {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error read file %q: %s", path, err.Error())
	}

	var result Root
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("error unmarshal file %q: %s", path, err.Error())
	}

	return result
}

func main() {
	ruBag := readBag("ru.json")
	enBag := readBag("en.json")

	cards := make(map[string]tts_entity.Card)
	debug := true

	uniqFill := make(map[string]int)

	// range via collections
	for _, col := range enBag.ObjectStates[0].ContainedObjects {
		//if debug {
		//	fmt.Println("collection:", col.Nickname)
		//}
		// range via decks inside collection
		for _, deck := range col.ContainedObjects {
			//if debug {
			//	fmt.Println(" -> deck:", deck.Nickname)
			//}
			// range via cards inside deck
			for _, card := range deck.ContainedObjects {
				//if debug {
				//	fmt.Println(" -> -> card:", card.Nickname)
				//}
				label := strings.ToLower(col.Nickname + "." + deck.Nickname + "." + card.Nickname)
				index, ok := uniqFill[label]
				if !ok {
					uniqFill[label] = 1
					cards[label] = card
				} else {
					uniqFill[label] += 1
					label += "." + strconv.Itoa(index)
					cards[label] = card
				}
			}
		}
	}

	uniqRead := make(map[string]int)

	// range via collections
	for colID, col := range ruBag.ObjectStates[0].ContainedObjects {
		if debug {
			fmt.Println("collection:", col.Nickname)
		}
		// range via decks inside collection
		for deckID, deck := range col.ContainedObjects {
			if debug {
				fmt.Println(" -> deck:", deck.Nickname)
			}
			// range via cards inside deck
			for cardID, card := range deck.ContainedObjects {
				if debug {
					fmt.Println(" -> -> card:", card.Nickname)
				}
				names := strings.Split(card.Nickname, "/")
				if len(names) != 2 && len(names) != 4 {
					log.Printf("invalid card name: %q", card.Nickname)
					continue
				}
				if len(names) == 4 {
					names = []string{
						strings.Join(names[:2], "/"),
						strings.Join(names[2:], "/"),
					}
				}

				names[0] = strings.TrimSpace(names[0])
				label := strings.ToLower(col.Nickname + "." + deck.Nickname + "." + names[0])
				index, ok := uniqRead[label]
				if !ok {
					uniqRead[label] = 1
				} else {
					uniqRead[label] += 1
					label += "." + strconv.Itoa(index)
				}
				if c, ok := cards[label]; ok {
					ruBag.ObjectStates[0].ContainedObjects[colID].ContainedObjects[deckID].ContainedObjects[cardID].States = map[string]tts_entity.Card{
						"2": c,
					}
					continue
				}

				names[1] = strings.TrimSpace(names[1])
				label = strings.ToLower(col.Nickname + "." + deck.Nickname + "." + names[1])
				index, ok = uniqRead[label]
				if !ok {
					uniqRead[label] = 1
				} else {
					uniqRead[label] += 1
					label += "." + strconv.Itoa(index)
				}
				if c, ok := cards[label]; ok {
					ruBag.ObjectStates[0].ContainedObjects[colID].ContainedObjects[deckID].ContainedObjects[cardID].States = map[string]tts_entity.Card{
						"2": c,
					}
					continue
				}
				log.Printf("state for card %q not found", card.Nickname)
			}
		}
	}

	data, err := json.MarshalIndent(ruBag, "", "\t")
	if err != nil {
		log.Fatal("error marshal:", err)
	}

	file, err := os.Create("join.json")
	if err != nil {
		log.Fatal("error create result file:", err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		log.Fatal("error write result data:", err)
	}
	err = file.Sync()
	if err != nil {
		log.Fatal("error sync:", err)
	}
}
