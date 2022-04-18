package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
	"tts_deck_build/internal/errors"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
func CleanTitle(in string) string {
	res := strings.ReplaceAll(in, " / ", "_")
	res = strings.ReplaceAll(res, "/", "_")
	res = strings.ReplaceAll(res, "!", "")
	res = strings.ReplaceAll(res, "'", "")
	res = strings.ReplaceAll(res, ".", "")
	return strings.ReplaceAll(res, " ", "_")
}
func GetFilenameFromUrl(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, filename := path.Split(u.Path)
	return filename
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
			IfErrorLog(resp.Body.Close())
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

func ToJson(data interface{}) (res []byte) {
	res, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
	}
	return
}
func RequestToObject(r io.ReadCloser, data interface{}) (e *errors.Error) {
	defer func() { IfErrorLog(r.Close()) }()
	err := json.NewDecoder(r).Decode(data)
	if err != nil {
		IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}
func IfErrorLog(err error) {
	if err != nil {
		log.Output(2, err.Error())
	}
}

func FileExist(path string) (isExist bool, e *errors.Error) {
	isExist, _, e = IsDir(path)
	return
}
func IsDir(path string) (isExist, isDir bool, e *errors.Error) {
	stat, err := os.Stat(path)
	if err == nil {
		isExist = true
	}

	if err != nil {
		if !os.IsNotExist(err) {
			// Some error
			IfErrorLog(err)
			e = errors.InternalError.AddMessage(err.Error())
		}
		// File is not exist
		return
	}

	if stat.IsDir() {
		isDir = true
	}
	return
}
func RemoveDir(path string) (e *errors.Error) {
	err := os.RemoveAll(path)
	if err != nil {
		IfErrorLog(err)
		e = errors.InternalError.AddMessage(err.Error())
	}
	return
}

func ResponseError(w http.ResponseWriter, e *errors.Error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.GetCode())
	if len(e.GetMessage()) > 0 {
		_, err := w.Write(ToJson(e))
		IfErrorLog(err)
	}
	return
}
