package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type (
	DownloadInfo struct {
		FilePath string
		FileName string
		URL      string
	}
)

// Download one file with passed name
func DownloadFile(filePath string, link string) {
	resp, err := http.Get(link)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Check which files already exists in cache
func checkCache() map[string]struct{} {
	res := make(map[string]struct{})

	files, err := ioutil.ReadDir(CachePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		res[file.Name()] = struct{}{}
	}
	return res
}

// Download all files
func DownloadFiles(pairs []DownloadInfo) {
	cache := checkCache()

	cachedFiles := 0
	downloadedFiles := 0
	for _, pair := range pairs {
		// Check if file already exists
		if _, ok := cache[pair.FileName]; ok {
			cachedFiles++
			continue
		}
		// Download if it's new file
		log.Println("File:", pair.FileName, "downloading...")
		DownloadFile(pair.FilePath, pair.URL)
		// Mark as downloaded
		cache[pair.FileName] = struct{}{}
		downloadedFiles++
	}
	log.Println("Total", cachedFiles, "images already exists in cache")
	log.Println("Total", downloadedFiles, "images were downloaded")
}
