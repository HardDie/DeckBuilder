package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/HardDie/fsentry/pkg/fsentry_types"
)

type Cards struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	Data      map[string]struct {
		ID          int                                   `json:"id"`
		Name        fsentry_types.QuotedString            `json:"name"`
		Description fsentry_types.QuotedString            `json:"description"`
		Image       fsentry_types.QuotedString            `json:"image"`
		Variables   map[string]fsentry_types.QuotedString `json:"variables"`
		Count       int                                   `json:"count"`
		CreatedAt   *time.Time                            `json:"createdAt"`
		UpdatedAt   *time.Time                            `json:"updatedAt"`
	} `json:"data"`
}

func readCards(path string) Cards {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	var ret Cards
	err = json.NewDecoder(file).Decode(&ret)
	if err != nil {
		log.Fatal(err)
	}

	return ret
}
func writeCards(path string, cards Cards) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(cards)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("example: app src/cards/.info.json dst/cards/.info.json")
		os.Exit(1)
	}

	// Read src file
	cardsSrc := readCards(os.Args[1])
	// Read dst file
	cardsDst := readCards(os.Args[2])

	// Copy variables
	for key, srcVal := range cardsSrc.Data {
		dstVal, ok := cardsDst.Data[key]
		if !ok {
			fmt.Println("card not exist:", key)
			continue
		}
		dstVal.Variables = srcVal.Variables
		cardsDst.Data[key] = dstVal
	}

	fmt.Println("dst cards")
	for key, val := range cardsDst.Data {
		fmt.Println(key, val.Name, val.Variables)
	}

	// Rewrite dst file
	writeCards(os.Args[2], cardsDst)
}
