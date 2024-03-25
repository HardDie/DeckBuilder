package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/HardDie/DeckBuilder/internal/tts_entity"
)

type Loop struct {
	uniqRead map[string]int
	ruBag    *Root
	cards    map[string]tts_entity.Card
}

func NewLoop(ruBag *Root, cards map[string]tts_entity.Card) *Loop {
	return &Loop{
		uniqRead: make(map[string]int),
		ruBag:    ruBag,
		cards:    cards,
	}
}

func (l *Loop) loopCollections(objects []CollectionBag) {
	for colID, col := range objects {
		if debug {
			fmt.Println("collection:", col.Nickname)
		}
		l.loopDeck(colID, col.Nickname, col.ContainedObjects)
	}
}

func (l *Loop) loopDeck(colID int, colNickname string, objects []tts_entity.TTSObject) {
root:
	for itemID, item := range objects {
		switch item.GetName() {
		case "Deck":
			deck := item.(tts_entity.DeckObject)
			l.loopCard(colID, itemID, colNickname, deck.Nickname, deck.ContainedObjects)
		case "Card":
			card := item.(tts_entity.Card)
			if debug {
				fmt.Println(" -> -> card:", card.Nickname)
			}

			names := getNames(card.Nickname)
			for _, name := range names {
				label := strings.ToLower(colNickname + "." + "" + "." + name)
				index, ok := l.uniqRead[label]
				if !ok {
					l.uniqRead[label] = 1
				} else {
					l.uniqRead[label] += 1
					label += "." + strconv.Itoa(index)
				}
				if c, ok := l.cards[label]; ok {
					cardSrc := l.ruBag.ObjectStates[0].ContainedObjects[colID].ContainedObjects[itemID].(tts_entity.Card)
					cardSrc.States = map[string]tts_entity.Card{
						"2": c,
					}
					l.ruBag.ObjectStates[0].ContainedObjects[colID].ContainedObjects[itemID] = cardSrc
					continue root
				}
			}

			log.Printf("state for card %q not found", card.Nickname)
		default:

		}
	}
}

func (l *Loop) loopCard(colID, deckID int, colNickname, deckNickname string, objects []tts_entity.Card) {
root:
	for cardID, card := range objects {
		if debug {
			fmt.Println(" -> -> card:", card.Nickname)
		}

		names := getNames(card.Nickname)
		for _, name := range names {
			label := strings.ToLower(colNickname + "." + deckNickname + "." + name)
			index, ok := l.uniqRead[label]
			if !ok {
				l.uniqRead[label] = 1
			} else {
				l.uniqRead[label] += 1
				label += "." + strconv.Itoa(index)
			}
			if c, ok := l.cards[label]; ok {
				deckSrc := l.ruBag.ObjectStates[0].ContainedObjects[colID].ContainedObjects[deckID].(tts_entity.DeckObject)
				deckSrc.ContainedObjects[cardID].States = map[string]tts_entity.Card{
					"2": c,
				}
				l.ruBag.ObjectStates[0].ContainedObjects[colID].ContainedObjects[deckID] = deckSrc
				continue root
			}
		}

		log.Printf("state for card %q not found", card.Nickname)
	}
}
func getNames(cardNickname string) []string {
	names := strings.Split(cardNickname, "/")
	if len(names) != 2 && len(names) != 4 {
		log.Printf("invalid card name: %q", cardNickname)
		return nil
	}
	if len(names) == 4 {
		names = []string{
			strings.Join(names[:2], "/"),
			strings.Join(names[2:], "/"),
		}
	}
	names[0] = strings.TrimSpace(names[0])
	names[1] = strings.TrimSpace(names[1])
	return names
}
