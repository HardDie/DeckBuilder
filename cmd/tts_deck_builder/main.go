//go:generate goversioninfo -icon=../../win_icon/icon.ico -64

// Package main TTS_deck_builder
//
// Entry point for the application.
//
// Terms Of Service:
//
//     Schemes: http
//     Host: localhost:5000
//     BasePath: /api
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
	"log"
	"net/http"

	"tts_deck_build/api"
	"tts_deck_build/internal/network"
)

func main() {
	// Setup logs
	log.SetFlags(log.Llongfile | log.Ltime)

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
