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
	"github.com/HardDie/DeckBuilder/internal/application"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
)

func main() {
	app, err := application.Get()
	if err != nil {
		logger.Error.Fatal(err.Error())
	}

	network.OpenBrowser("http://127.0.0.1:5000")

	err = app.Run()
	if err != nil {
		logger.Error.Fatal(err.Error())
	}
}
