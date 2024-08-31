package progress

import (
	"sync"

	entitiesStatus "github.com/HardDie/DeckBuilder/internal/entities/status"
	"github.com/HardDie/DeckBuilder/internal/logger"
)

var progressSingleton *progress

const (
	StatusEmpty      string = "empty"
	StatusInProgress        = "in_progress"
	StatusDone              = "done"
	StatusError             = "error"
)

type progress struct {
	Type     string
	Message  string
	Progress float32
	Status   string
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
	p.Status = StatusEmpty
}

func (p *progress) SetType(value string) {
	p.m.Lock()
	defer p.m.Unlock()
	p.Type = value

	logger.Debug.Printf("Progress(%s): [%s] %s - %0.2f\n", p.Status, p.Type, p.Message, p.Progress)
}

func (p *progress) SetMessage(value string) {
	p.m.Lock()
	defer p.m.Unlock()
	p.Message = value

	logger.Debug.Printf("Progress(%s): [%s] %s - %0.2f\n", p.Status, p.Type, p.Message, p.Progress)
}

func (p *progress) SetProgress(value float32) {
	p.m.Lock()
	defer p.m.Unlock()
	p.Progress = value

	logger.Debug.Printf("Progress(%s): [%s] %s - %0.2f\n", p.Status, p.Type, p.Message, p.Progress)
}

func (p *progress) SetStatus(value string) {
	p.m.Lock()
	defer p.m.Unlock()
	p.Status = value

	logger.Debug.Printf("Progress(%s): [%s] %s - %0.2f\n", p.Status, p.Type, p.Message, p.Progress)
}

func (p *progress) GetStatus() entitiesStatus.Status {
	return entitiesStatus.Status{
		Type:     p.Type,
		Message:  p.Message,
		Progress: p.Progress,
		Status:   p.Status,
	}
}
