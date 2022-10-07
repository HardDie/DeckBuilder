package web

import (
	"embed"
	"net/http"

	"github.com/HardDie/DeckBuilder/internal/logger"
)

var (
	//go:embed dist
	res embed.FS

	//go:embed swagger.json
	swagger []byte
)

func ServeSwagger(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(swagger)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}

func GetWeb() embed.FS {
	return res
}
