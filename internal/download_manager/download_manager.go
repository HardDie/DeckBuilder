package download_manager

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"tts_deck_build/internal/config"
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

func (m *DownloadManager) AddFile(localurl, filename string) {
	if _, ok := m.unique[filename]; ok {
		if config.GetConfig().Debug {
			log.Println("File already exist in queue:", filename)
		}
		return
	}

	_, err := url.Parse(localurl)
	if err != nil {
		log.Fatalf("Bad URL: %q %v", localurl, err.Error())
	}

	m.files = append(m.files, fileInfo{
		Filename: filename,
		URL:      localurl,
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
		u, _ := url.Parse(file.URL)
		if u.Scheme == "" {
			m.copy(filepath.Join(config.GetConfig().SourceDir, file.URL), file.Filename)
		} else {
			m.download(file.URL, file.Filename)
		}
		downloadedFiles++
	}
	log.Println("Total", cachedFiles, "images already exists in cache")
	log.Println("Total", downloadedFiles, "images were downloaded")
}

func (m *DownloadManager) copy(path, filename string) {
	// Check source file
	sourceFileStat, err := os.Stat(path)
	if err != nil {
		log.Fatal("Can't open file:", err.Error())
	}
	// Check this is regular file
	if !sourceFileStat.Mode().IsRegular() {
		log.Fatalf("%s is not a regular file", path)
	}
	// Open source file
	source, err := os.Open(path)
	if err != nil {
		log.Fatal("Can't open file:", err.Error())
	}
	defer source.Close()

	// Create destination file
	destination, err := os.Create(filepath.Join(m.cachePath, filename))
	if err != nil {
		log.Fatal("Can't create dest file:", err.Error())
	}
	defer destination.Close()

	// Copy data from source to destination file
	_, err = io.Copy(destination, source)
	if err != nil {
		log.Fatal("Can't copy file:", err.Error())
	}
}
func (m *DownloadManager) download(url, filename string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath.Join(m.cachePath, filename))
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
	files, err := ioutil.ReadDir(config.GetConfig().CachePath)
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
