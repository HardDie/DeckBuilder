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
func ResponseError(w http.ResponseWriter, e error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch val := e.(type) {
	case errors.Err:
		if val.GetCode() > 0 {
			w.WriteHeader(val.GetCode())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if len(val.GetMessage()) > 0 {
			_, err := w.Write(toJson(e))
			errors.IfErrorLog(err)
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		msg := "unhandled error: " + e.Error()
		log.Print(msg)
		_, err := w.Write([]byte(msg))
		errors.IfErrorLog(err)
	}
	return
}
func Response(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write(toJson(data))
	errors.IfErrorLog(err)
	return
}
