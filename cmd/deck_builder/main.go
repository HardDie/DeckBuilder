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
//
// swagger:meta
package main

import (
	"github.com/HardDie/DeckBuilder/internal/application"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/pkg/version"
)

func main() {
	app, err := application.Get(version.String())
	if err != nil {
		logger.Error.Fatal(err.Error())
	}

	err = app.Run()
	if err != nil {
		logger.Error.Fatal(err.Error())
	}
}
