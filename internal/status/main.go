package status

import (
	"log"
	"sync"
)

var statusSingleton *status

type status struct {
	Type     string
	Message  string
	Progress float32
	m        sync.Mutex
}

func GetStatus() *status {
	if statusSingleton != nil {
		return statusSingleton
	}
	statusSingleton = &status{}
	statusSingleton.Flush()
	return statusSingleton
}

func (s *status) Flush() {
	s.m.Lock()
	defer s.m.Unlock()
	s.Type = "No process"
	s.Message = ""
	s.Progress = 0
}

func (s *status) SetType(value string) {
	s.m.Lock()
	defer s.m.Unlock()
	s.Type = value

	log.Printf("Progress: [%s] %s - %0.2f\n", s.Type, s.Message, s.Progress)
}

func (s *status) SetMessage(value string) {
	s.m.Lock()
	defer s.m.Unlock()
	s.Message = value

	log.Printf("Progress: [%s] %s - %0.2f\n", s.Type, s.Message, s.Progress)
}

func (s *status) SetProgress(value float32) {
	s.m.Lock()
	defer s.m.Unlock()
	s.Progress = value

	log.Printf("Progress: [%s] %s - %0.2f\n", s.Type, s.Message, s.Progress)
}
