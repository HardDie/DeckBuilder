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

const (
	fillDebug = !true
	debug     = !true
)

func main() {
	ruBag := RawToValid(readBag("ru.json"))
	enBag := RawToValid(readBag("en.json"))

	cards := make(map[string]tts_entity.Card)

	uniqFill := make(map[string]int)

	// range via collections
	for _, col := range enBag.ObjectStates[0].ContainedObjects {
		if fillDebug {
			fmt.Println("collection:", col.Nickname)
		}
		// range via items inside collection
		for _, item := range col.ContainedObjects {
			switch item.GetName() {
			case "Deck":
				deck := item.(tts_entity.DeckObject)
				if fillDebug {
					fmt.Println(" -> deck:", deck.Nickname)
				}
				// range via cards inside deck
				for _, card := range deck.ContainedObjects {
					if fillDebug {
						fmt.Println(" -> -> card:", card.Nickname)
					}
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
			case "Card":
				card := item.(tts_entity.Card)
				if fillDebug {
					fmt.Println(" -> -> card:", card.Nickname)
				}
				label := strings.ToLower(col.Nickname + "." + "" + "." + card.Nickname)
				index, ok := uniqFill[label]
				if !ok {
					uniqFill[label] = 1
					cards[label] = card
				} else {
					uniqFill[label] += 1
					label += "." + strconv.Itoa(index)
					cards[label] = card
				}
			default:
				log.Fatal("unknown item type:", item.GetName())
			}
		}
	}

	NewLoop(&ruBag, cards).loopCollections(ruBag.ObjectStates[0].ContainedObjects)

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
