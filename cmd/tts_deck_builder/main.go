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
	"tts_deck_build/internal/application"
	"tts_deck_build/internal/logger"
	"tts_deck_build/internal/network"
)

func main() {
	app, err := application.Get()
	if err != nil {
		logger.Error.Fatal(err.Error())
	}

	if false {
		network.OpenBrowser("http://localhost:5000")
	}

	err = app.Run()
	if err != nil {
		logger.Error.Fatal(err.Error())
	}
}
