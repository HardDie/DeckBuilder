//go:generate goversioninfo -icon=../../deployment/win_icon.ico -64

// Package main DeckBuilder
//
// Entry point for the application.
//
// Terms Of Service:
//
//	Schemes: http
//	Host: localhost:5000
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//	- binary
//
// swagger:meta
package main

import (
	"flag"

	"github.com/HardDie/DeckBuilder/internal/application"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
)

func main() {
	// If the flag is set, run the game in debug mode.
	// - Do not request the url and don't open the browser
	// - Do not close the application when /system/quit is requested
	debugFlag := flag.Bool("debug", false, "")
	flag.Parse()

	app, err := application.Get(*debugFlag)
	if err != nil {
		logger.Error.Fatal(err.Error())
	}

	if !*debugFlag {
		network.OpenBrowser("http://127.0.0.1:5000")
	}

	err = app.Run()
	if err != nil {
		logger.Error.Fatal(err.Error())
	}
}
