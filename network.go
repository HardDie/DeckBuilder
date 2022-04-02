package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type fileInfo struct {
	Filename string
	URL      string
}

type DownloadManager struct {
	unique    map[string]struct{}
	cache     map[string]struct{}
	files     []fileInfo
	cachePath string
}

func NewDownloadManager(cachePath string) *DownloadManager {
	return &DownloadManager{
		unique:    make(map[string]struct{}),
		cache:     make(map[string]struct{}),
		cachePath: cachePath,
	}
}

func (m *DownloadManager) AddFile(url, filename string) {
	if _, ok := m.unique[filename]; ok {
		if GetConfig().Debug {
			log.Println("File already exist in queue:", filename)
		}
		return
	}

	m.files = append(m.files, fileInfo{
		Filename: filename,
		URL:      url,
	})
	m.unique[filename] = struct{}{}
}
func (m *DownloadManager) Download() {
	m.checkCache()

	cachedFiles := 0
	downloadedFiles := 0
	for _, file := range m.files {
		// Check if file already exists
		if _, ok := m.cache[file.Filename]; ok {
			cachedFiles++
			continue
		}
		// Download if it's new file
		log.Println("File:", file.Filename, "downloading...")
		m.download(file.URL, file.Filename)
		downloadedFiles++
	}
	log.Println("Total", cachedFiles, "images already exists in cache")
	log.Println("Total", downloadedFiles, "images were downloaded")
}

func (m *DownloadManager) download(url, filename string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	out, err := os.Create(m.cachePath + filename)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	m.cache[filename] = struct{}{}
}
func (m *DownloadManager) checkCache() {
	files, err := ioutil.ReadDir(GetConfig().CachePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		m.cache[file.Name()] = struct{}{}
	}
}
