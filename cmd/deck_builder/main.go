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
	"runtime/debug"

	"github.com/HardDie/DeckBuilder/internal/application"
	"github.com/HardDie/DeckBuilder/internal/logger"
	"github.com/HardDie/DeckBuilder/internal/network"
)

var (
	Version        = ""
	BackendCommit  = ""
	FrontendCommit = ""
)

func main() {
	// If the flag is set, run the game in debug mode.
	// - Do not request the url and don't open the browser
	// - Do not close the application when /system/quit is requested
	debugFlag := flag.Bool("debug", false, "")
	flag.Parse()

	if info, available := debug.ReadBuildInfo(); available {
		switch info.Main.Version {
		case "", "(devel)":
			// skip
		default:
			// In case we installed the application as "go install ..." from github
			Version = info.Main.Version
		}
	}

	var version string
	if BackendCommit != "" {
		// If the application was built using the deployment script
		version = "Backend: " + BackendCommit + ", Frontend: " + FrontendCommit
	} else if Version != "" {
		// If the application was installed as a "go install ..."
		version = Version
	} else {
		// Bad case
		version = "unknown"
	}

	app, err := application.Get(*debugFlag, version)
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
