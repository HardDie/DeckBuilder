package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Read configurations, download images, build deck image files
func GenerateDeckImages() {
	// Read all decks
	listOfDecks := Crawl(GetConfig().SourceDir)

	dm := NewDownloadManager(GetConfig().CachePath)
	// Fill download list
	for _, decks := range listOfDecks {
		for _, deck := range decks {
			PutDeckToDownloadManager(deck, dm)
		}
	}
	// Download all images
	dm.Download()

	// Build
	db := NewDeckBuilder()
	for _, decks := range listOfDecks {
		for _, deck := range decks {
			PutDeckToDeckBuilder(deck, db)
		}
	}

	// Generate images
	images := db.DrawDecks()

	// Write all created files
	data, _ := json.MarshalIndent(images, "", "	")
	err := ioutil.WriteFile(filepath.Join(GetConfig().ResultDir, "images.json"), data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Read configurations, generate TTS json object with description
func GenerateDeckObject() {
	// Read all decks
	listOfDecks := Crawl(GetConfig().SourceDir)

	// Build
	db := NewDeckBuilder()
	for _, decks := range listOfDecks {
		for _, deck := range decks {
			PutDeckToDeckBuilder(deck, db)
		}
	}

	// Generate TTS object
	res := db.GenerateTTSDeck()

	// Write deck json to file
	err := ioutil.WriteFile(filepath.Join(GetConfig().ResultDir, "deck.json"), res, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	err := exec.Command(cmd, args...).Start()
	if err != nil {
		log.Fatal("Can't run browser")
	}
}

func WebServer() {
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	log.Println("Listening on :5000...")

	go func() {
		for {
			resp, err := http.Get("http://localhost:5000")
			if err != nil {
				log.Println("Failed:", err)
				continue
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Println("Not OK:", resp.StatusCode)
				continue
			}

			// Reached this point: server is up and running!
			break
		}
		openBrowser("http://localhost:5000")
	}()

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	return
}

func main() {
	// Setup logs
	log.SetFlags(log.Lshortfile | log.Ltime)

	// Setup run flags
	genImgMode := flag.Bool("generate_image", false, "Run process of generating deck images")
	genDeckMode := flag.Bool("generate_object", false, "Run process of generating json deck object")
	helpMode := flag.Bool("help", false, "Show help")
	flag.Parse()

	switch {
	case *genImgMode:
		GenerateDeckImages()
	case *genDeckMode:
		GenerateDeckObject()
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
