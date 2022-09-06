//go:generate goversioninfo -icon=../../win_icon/icon.ico -64

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
	generator "tts_deck_build/internal/generator_old"
	"tts_deck_build/internal/network"
)

func WebServer() {
	log.Println("Listening on :5000...")

	if false {
		network.OpenBrowser("http://localhost:5000")
	}

	http.Handle("/", api.GetRoutes())
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	// Setup logs
	log.SetFlags(log.Llongfile | log.Ltime)

	// Setup run flags
	genDeckMode := flag.Bool("generate_object", false, "Run process of generating json deck object")
	gameName := flag.String("game", "", "Title of game for generator")
	helpMode := flag.Bool("help", false, "Show help")
	flag.Parse()

	switch {
	case *genDeckMode:
		generator.GenerateDeckObject(*gameName)
	case *helpMode:
		fmt.Println("How to use:")
		fmt.Println("1. Build images from ${sourceDir}/*.json descriptions (-generate_image -game gameID)")
		fmt.Println("2. Upload images on some hosting (steam cloud)")
		fmt.Println("3. Write URL for each image in ${resultDir}/images.json file")
		fmt.Println("4. Build deck object ${resultDir}/deck.json (-generate_object -game gameID)")
		fmt.Println("5. Put deck object into \"Tabletop Simulator/Saves/Saved Objects\" folder")
		fmt.Println()
		fmt.Println("Choose one of the mode:")
		flag.PrintDefaults()
	default:
		WebServer()
	}
}
