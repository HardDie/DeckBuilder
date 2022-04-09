package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
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

func ToJson(data interface{}) (res string) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
	}
	res = string(jsonData)
	return
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
			resp.Body.Close()
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
