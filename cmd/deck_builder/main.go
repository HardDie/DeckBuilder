// Image and TTS HTTP server. No window.
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
