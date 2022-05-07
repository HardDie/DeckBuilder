package network

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"tts_deck_build/internal/errors"
)

type JsonResponse struct {
	// Body
	Data interface{} `json:"data,omitempty"`
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
		log.Fatal("Can't run browser")
	}
}
func OpenBrowser(url string) {
	go func() {
		for {
			time.Sleep(time.Millisecond)
			resp, err := http.Get(url)
			if err != nil {
				log.Println("Failed:", err)
				continue
			}
			errors.IfErrorLog(resp.Body.Close())
			if resp.StatusCode != http.StatusOK {
				log.Println("Not OK:", resp.StatusCode)
				continue
			}

			// Reached this point: server is up and running!
			break
		}
		openBrowser(url)
	}()
}

func toJson(data interface{}) (res []byte) {
	res, err := json.Marshal(data)
	if err != nil {
		errors.IfErrorLog(err)
	}
	return
}
func RequestToObject(r io.ReadCloser, data interface{}) (e error) {
	defer func() { errors.IfErrorLog(r.Close()) }()
	err := json.NewDecoder(r).Decode(data)
	if err != nil {
		errors.IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}
func response(w http.ResponseWriter, httpCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)
	_, err := w.Write(toJson(data))
	errors.IfErrorLog(err)
	return err
}
func ResponseError(w http.ResponseWriter, e error) {
	resp := JsonResponse{
		Error: e,
	}

	httpCode := http.StatusInternalServerError
	if val, ok := e.(*errors.Err); ok {
		if val.GetCode() > 0 {
			httpCode = val.GetCode()
		}
	} else {
		log.Println("unhandled error: " + e.Error())
	}

	_ = response(w, httpCode, resp)
	return
}
func Response(w http.ResponseWriter, data interface{}) {
	resp := JsonResponse{
		Data: data,
	}

	_ = response(w, http.StatusOK, resp)
	return
}
