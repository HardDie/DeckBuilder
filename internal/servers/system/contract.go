package system

import "net/http"

type System interface {
	QuitHandler(w http.ResponseWriter, r *http.Request)
	GetSettingsHandler(w http.ResponseWriter, r *http.Request)
	UpdateSettingsHandler(w http.ResponseWriter, r *http.Request)
	StatusHandler(w http.ResponseWriter, r *http.Request)
	GetVersionHandler(w http.ResponseWriter, r *http.Request)
	StopQuit()
}
