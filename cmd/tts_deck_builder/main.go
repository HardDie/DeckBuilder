// Package main TTS_deck_builder
//
// Entry point for the application.
//
// Terms Of Service:
//
//     Schemes: http
//     Host: localhost:5000
//     BasePath: /
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     - binary
//
// swagger:meta
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"tts_deck_build/api"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/generator"
	"tts_deck_build/internal/network"
)

func WebServer() {
	log.Println("Listening on :5000...")

	network.OpenBrowser("http://localhost:5000")

	http.Handle("/", api.GetRoutes())
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	return
}

func createDirIfNotExists(folder string) {
	isExists, e := fs.IsFolderExist(folder)
	if e != nil {
		log.Fatal(e.Error())
	}
	if isExists {
		return
	}
	e = fs.CreateFolder(folder)
	if e != nil {
		log.Fatal(e.Error())
	}
}
func setup() {
	createDirIfNotExists(config.GetConfig().Data)
	createDirIfNotExists(config.GetConfig().Games())
}

func main() {
	// Setup logs
	log.SetFlags(log.Llongfile | log.Ltime)

	setup()

	// Setup run flags
	genImgMode := flag.Bool("generate_image", false, "Run process of generating deck images")
	genDeckMode := flag.Bool("generate_object", false, "Run process of generating json deck object")
	helpMode := flag.Bool("help", false, "Show help")
	flag.Parse()

	switch {
	case *genImgMode:
		generator.GenerateDeckImages()
	case *genDeckMode:
		generator.GenerateDeckObject()
	case *helpMode:
		fmt.Println("How to use:")
		fmt.Println("1. Build images from ${sourceDir}/*.json descriptions (-generate_image)")
		fmt.Println("2. Upload images on some hosting (steam cloud)")
		fmt.Println("3. Write URL for each image in ${resultDir}/images.json file")
		fmt.Println("4. Build deck object ${resultDir}/deck.json (-generate_object)")
		fmt.Println("5. Put deck object into \"Tabletop Simulator/Saves/Saved Objects\" folder")
		fmt.Println()
		fmt.Println("Choose one of the mode:")
		flag.PrintDefaults()
	default:
		WebServer()
	}
}
