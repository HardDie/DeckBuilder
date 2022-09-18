package web

import (
	"embed"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"

	"tts_deck_build/internal/logger"
)

var (
	//go:embed web
	res embed.FS

	pages = map[string]string{}
)

func servePages(w http.ResponseWriter, r *http.Request) {
	page, ok := pages[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	file, err := res.ReadFile(page)
	if err != nil {
		logger.Error.Printf("page %s not found in pages cache...", r.RequestURI)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch filepath.Ext(page) {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "text/javascript")
	case ".json":
		w.Header().Set("Content-Type", "application/json")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	}
	_, err = w.Write(file)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

func registerFiles(dirName string) {
	files, _ := res.ReadDir(dirName)
	for _, file := range files {
		if file.IsDir() {
			registerFiles(filepath.Join(dirName, file.Name()))
		} else {
			fullName := filepath.Join(dirName, file.Name())
			src := strings.TrimPrefix(fullName, "web")
			if src == "/index.html" {
				pages["/"] = fullName
				continue
			}
			pages[src] = fullName
		}
	}
}

func Init(route *mux.Router) {
	registerFiles("web")
	for page := range pages {
		route.HandleFunc(page, servePages)
	}

	// DEVELOP PURPOSE ONLY
	redocHandler := middleware.Redoc(middleware.RedocOpts{SpecURL: "/swagger.json"}, nil)
	route.Handle("/docs", redocHandler)
	// DEVELOP PURPOSE ONLY
}
