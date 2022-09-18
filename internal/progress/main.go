package progress

import (
	"sync"

	"tts_deck_build/internal/logger"
)

var progressSingleton *progress

type progress struct {
	Type     string
	Message  string
	Progress float32
	m        sync.Mutex
}

func GetProgress() *progress {
	if progressSingleton != nil {
		return progressSingleton
	}
	progressSingleton = &progress{}
	progressSingleton.Flush()
	return progressSingleton
}

func (p *progress) Flush() {
	p.m.Lock()
	defer p.m.Unlock()
	p.Type = "No process"
	p.Message = ""
	p.Progress = 0
}

func (p *progress) SetType(value string) {
	p.m.Lock()
	defer p.m.Unlock()
	p.Type = value

	logger.Debug.Printf("Progress: [%s] %s - %0.2f\n", p.Type, p.Message, p.Progress)
}

func (p *progress) SetMessage(value string) {
	p.m.Lock()
	defer p.m.Unlock()
	p.Message = value

	logger.Debug.Printf("Progress: [%s] %s - %0.2f\n", p.Type, p.Message, p.Progress)
}

func (p *progress) SetProgress(value float32) {
	p.m.Lock()
	defer p.m.Unlock()
	p.Progress = value

	logger.Debug.Printf("Progress: [%s] %s - %0.2f\n", p.Type, p.Message, p.Progress)
}
