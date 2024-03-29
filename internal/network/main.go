package network

import (
	"encoding/json"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/HardDie/DeckBuilder/internal/errors"
	"github.com/HardDie/DeckBuilder/internal/fs"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

type Meta struct {
	Total      int `json:"total"`
	CardsTotal int `json:"cardsTotal,omitempty"`
	//Limit int `json:"limit"`
	//Page  int `json:"page"`
}
type JSONResponse struct {
	// Body
	Data interface{} `json:"data,omitempty"`
	// Meta
	Meta *Meta `json:"meta,omitempty"`
	// Error information
	Error interface{} `json:"error,omitempty"`
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
		logger.Error.Fatal("Can't run browser")
	}
}
func OpenBrowser(url string) {
	go func() {
		for {
			time.Sleep(time.Millisecond)
			resp, err := http.Get(url)
			if err != nil {
				logger.Info.Println("Failed:", err)
				continue
			}
			errors.IfErrorLog(resp.Body.Close())
			if resp.StatusCode != http.StatusOK {
				logger.Info.Println("Not OK:", resp.StatusCode)
				continue
			}

			// Reached this point: server is up and running!
			break
		}
		openBrowser(url)
	}()
}

func RequestToObject(r io.ReadCloser, data interface{}) (e error) {
	defer func() { errors.IfErrorLog(r.Close()) }()
	err := json.NewDecoder(r).Decode(data)
	if err != nil {
		if err == io.EOF {
			// If the body of the request is empty, it is not an error
			return nil
		}
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}
func response(w http.ResponseWriter, httpCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)
	return fs.JsonToWriter(w, data)
}
func ResponseError(w http.ResponseWriter, e error) {
	resp := JSONResponse{
		Error: e,
	}

	httpCode := http.StatusInternalServerError
	if val, ok := e.(*errors.Err); ok {
		if val.GetCode() > 0 {
			httpCode = val.GetCode()
		}
	} else {
		logger.Warn.Println("unhandled error: " + e.Error())
	}

	_ = response(w, httpCode, resp)
}
func Response(w http.ResponseWriter, data interface{}) {
	resp := JSONResponse{
		Data: data,
	}

	_ = response(w, http.StatusOK, resp)
}
func ResponseWithMeta(w http.ResponseWriter, data interface{}, meta *Meta) {
	resp := JSONResponse{
		Data: data,
		Meta: meta,
	}

	_ = response(w, http.StatusOK, resp)
}
